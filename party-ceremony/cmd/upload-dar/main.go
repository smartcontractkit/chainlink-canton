// upload-dar uploads one or more DAR files to a Canton participant via Admin PackageService.
//
// Usage:
//
//	go run ./cmd/upload-dar \
//	  -config config/staging/nodes/cv0.participant-config.json \
//	  -dir ../../contracts/dars/current
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

type darList []string

func (d *darList) String() string { return fmt.Sprint([]string(*d)) }

func (d *darList) Set(v string) error {
	*d = append(*d, v)
	return nil
}

func main() {
	configPath := flag.String("config", "", "path to participant-config.json")
	dir := flag.String("dir", "", "upload all .dar files in this directory (one request per file)")
	var dars darList
	flag.Var(&dars, "dar", "path to a DAR file (repeatable)")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "-config is required")
		os.Exit(2)
	}

	paths := append([]string{}, dars...)
	if *dir != "" {
		entries, err := filepath.Glob(filepath.Join(*dir, "*.dar"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "glob %q: %v\n", *dir, err)
			os.Exit(1)
		}
		if len(entries) == 0 {
			fmt.Fprintf(os.Stderr, "no .dar files in %q\n", *dir)
			os.Exit(1)
		}
		sort.Strings(entries)
		paths = append(paths, entries...)
	}
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "provide -dar and/or -dir")
		os.Exit(2)
	}

	cfg, err := client.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	conn, err := client.Dial(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect admin API: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	admin := client.NewGRPCClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
			os.Exit(1)
		}
		darID, err := admin.UploadDar(ctx, data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "upload %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("uploaded %s -> %s\n", filepath.Base(path), darID)
	}
}
