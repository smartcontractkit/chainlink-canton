package utilitydars

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	resolveDirMu       sync.RWMutex
	resolveDirOverride string
)

// SetResolveDir overrides the directory used by ResolveDirectory (e.g. ephemeral tmp during binding generation).
// Pass an empty string to clear the override.
func SetResolveDir(dir string) {
	resolveDirMu.Lock()
	defer resolveDirMu.Unlock()
	resolveDirOverride = dir
}

// ResolveDirectory returns the active utility DAR directory (override or process cache).
func ResolveDirectory() (string, error) {
	resolveDirMu.RLock()
	override := resolveDirOverride
	resolveDirMu.RUnlock()
	if override != "" {
		return override, nil
	}

	return cacheDir()
}

// FetchOptions configures a utility bundle fetch.
type FetchOptions struct {
	Dir                string
	PackageIDs         map[string]string // package name -> expected main package ID
	SkipPackageIDCheck bool
}

// Fetch downloads the pinned bundle, verifies its SHA256, and extracts selected DARs into Dir.
func Fetch(ctx context.Context, opts FetchOptions) error {
	m, err := currentManifest()
	if err != nil {
		return err
	}
	if opts.Dir == "" {
		return fmt.Errorf("utility fetch: Dir is required")
	}

	if err := os.MkdirAll(opts.Dir, 0o755); err != nil {
		return fmt.Errorf("create utility DAR dir %s: %w", opts.Dir, err)
	}

	tarball, err := downloadBundle(ctx, m.URL)
	if err != nil {
		return err
	}
	if err := verifySHA256(tarball, m.SHA256); err != nil {
		return err
	}
	if err := extractPackages(tarball, m, opts.Dir); err != nil {
		return err
	}

	if opts.SkipPackageIDCheck {
		return nil
	}

	for name, semver := range m.Packages {
		path := ResolveDarPath(name, semver, opts.Dir)
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read extracted DAR %s: %w", path, err)
		}
		expectedID, ok := opts.PackageIDs[name]
		if !ok {
			return fmt.Errorf("utility fetch: missing expected package ID for %q", name)
		}
		if err := VerifyPackageID(data, expectedID); err != nil {
			return fmt.Errorf("verify %s: %w", name, err)
		}
	}

	return nil
}

// EnsureCache populates the process cache if required DARs are missing or invalid.
func EnsureCache(ctx context.Context, packageIDs map[string]string) error {
	dir, err := cacheDir()
	if err != nil {
		return err
	}

	m, err := currentManifest()
	if err != nil {
		return err
	}

	if cacheValid(dir, m, packageIDs) {
		return nil
	}

	return Fetch(ctx, FetchOptions{
		Dir:        dir,
		PackageIDs: packageIDs,
	})
}

// FetchToTemp downloads and extracts utility DARs into a temporary directory.
// The caller must invoke cleanup when finished.
func FetchToTemp(ctx context.Context, packageIDs map[string]string) (dir string, cleanup func(), err error) {
	tmpDir, err := os.MkdirTemp("", "chainlink-canton-utility-dars-*")
	if err != nil {
		return "", nil, fmt.Errorf("create utility temp dir: %w", err)
	}

	cleanup = func() {
		_ = os.RemoveAll(tmpDir)
	}

	if err := Fetch(ctx, FetchOptions{
		Dir:        tmpDir,
		PackageIDs: packageIDs,
	}); err != nil {
		cleanup()
		return "", nil, err
	}

	return tmpDir, cleanup, nil
}

func cacheDir() (string, error) {
	m, err := currentManifest()
	if err != nil {
		return "", err
	}

	base, err := os.UserCacheDir()
	if err != nil {
		base = os.TempDir()
	}

	return filepath.Join(base, "chainlink-canton", fmt.Sprintf("utility-dars-%s", m.BundleVersion)), nil
}

func cacheValid(dir string, m *Manifest, packageIDs map[string]string) bool {
	for name, semver := range m.Packages {
		path := ResolveDarPath(name, semver, dir)
		data, err := os.ReadFile(path)
		if err != nil {
			return false
		}
		expectedID, ok := packageIDs[name]
		if !ok {
			return false
		}
		if err := VerifyPackageID(data, expectedID); err != nil {
			return false
		}
	}

	return true
}

func downloadBundle(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create utility bundle request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download utility bundle: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download utility bundle: HTTP %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read utility bundle: %w", err)
	}

	return data, nil
}

func verifySHA256(data []byte, expectedHex string) error {
	sum := sha256.Sum256(data)
	actual := hex.EncodeToString(sum[:])
	if !strings.EqualFold(actual, expectedHex) {
		return fmt.Errorf("utility bundle sha256 mismatch: got %s, want %s", actual, expectedHex)
	}

	return nil
}

func extractPackages(tarball []byte, m *Manifest, dir string) error {
	gz, err := gzip.NewReader(bytes.NewReader(tarball))
	if err != nil {
		return fmt.Errorf("open utility bundle gzip: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	want := make(map[string]string, len(m.Packages))
	for name, semver := range m.Packages {
		want[fmt.Sprintf("%s-%s.dar", name, semver)] = ResolveDarPath(name, semver, dir)
	}

	found := make(map[string]struct{}, len(want))
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read utility bundle tar: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		base := filepath.Base(hdr.Name)
		dest, ok := want[base]
		if !ok {
			continue
		}

		out, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("create %s: %w", dest, err)
		}

		if _, err := io.CopyN(out, tr, hdr.Size); err != nil {
			_ = out.Close()
			return fmt.Errorf("extract %s: %w", dest, err)
		}
		if err := out.Close(); err != nil {
			return fmt.Errorf("close %s: %w", dest, err)
		}

		found[base] = struct{}{}
	}

	for name := range want {
		if _, ok := found[name]; !ok {
			return fmt.Errorf("utility bundle missing %s", name)
		}
	}

	return nil
}
