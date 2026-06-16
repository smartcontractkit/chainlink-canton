package scriptenv

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Load loads the given .env file if it exists. Existing process env vars win.
func Load(envRelPath string) (string, error) {
	path, err := findEnv(envRelPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	return path, loadFile(path)
}

func loadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("%s:%d: expected KEY=VALUE", path, lineNo)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return fmt.Errorf("%s:%d: empty env key", path, lineNo)
		}

		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			unquoted, err := strconv.Unquote(value)
			if err != nil {
				return fmt.Errorf("%s:%d: invalid quoted value for %s: %w", path, lineNo, key, err)
			}
			value = unquoted
		} else if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
			value = value[1 : len(value)-1]
		}

		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("%s:%d: set %s: %w", path, lineNo, key, err)
			}
		}
	}

	return scanner.Err()
}

func String(defaultValue string, keys ...string) string {
	if value := First(keys...); value != "" {
		return value
	}
	return defaultValue
}

func Uint64(defaultValue uint64, keys ...string) (uint64, error) {
	value := First(keys...)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s as uint64: %w", keys[0], err)
	}
	return parsed, nil
}

func Int64(defaultValue int64, keys ...string) (int64, error) {
	value := First(keys...)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s as int64: %w", keys[0], err)
	}
	return parsed, nil
}

func Duration(defaultValue time.Duration, keys ...string) (time.Duration, error) {
	value := First(keys...)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s as duration: %w", keys[0], err)
	}
	return parsed, nil
}

func First(keys ...string) string {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func findEnv(envRelPath string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for dir := wd; ; dir = filepath.Dir(dir) {
		path := filepath.Join(dir, envRelPath)
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			return path, nil
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}

	return "", os.ErrNotExist
}
