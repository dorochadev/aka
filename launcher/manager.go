package launcher

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetLauncherDir() string {
	if dir := os.Getenv("AKA_BIN_DIR"); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("/tmp", "aka-launchers")
	}
	return filepath.Join(home, "bin")
}

func EnsureLauncherDir() error {
	dir := GetLauncherDir()
	return os.MkdirAll(dir, 0755)
}

func Exists(name string) bool {
	path := filepath.Join(GetLauncherDir(), name)
	_, err := os.Stat(path)
	return err == nil
}

func Create(name, appName string) error {
	if err := EnsureLauncherDir(); err != nil {
		return fmt.Errorf("failed to create launcher directory: %w", err)
	}

	path := filepath.Join(GetLauncherDir(), name)
	script := GenerateScript(appName)

	if err := os.WriteFile(path, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to write launcher file: %w", err)
	}

	return nil
}

func Remove(name string) error {
	path := filepath.Join(GetLauncherDir(), name)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove launcher: %w", err)
	}
	return nil
}

func Rename(oldName, newName string) error {
	oldPath := filepath.Join(GetLauncherDir(), oldName)
	newPath := filepath.Join(GetLauncherDir(), newName)

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename launcher: %w", err)
	}

	return nil
}

type LauncherInfo struct {
	Name   string
	Target string
}

func List() ([]LauncherInfo, error) {
	dir := GetLauncherDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []LauncherInfo{}, nil
		}
		return nil, fmt.Errorf("failed to read launcher directory: %w", err)
	}

	var launchers []LauncherInfo
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		target, err := extractTarget(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}

		launchers = append(launchers, LauncherInfo{
			Name:   entry.Name(),
			Target: target,
		})
	}

	return launchers, nil
}

// extractTarget reads a launcher file and extracts the target application name
func extractTarget(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if runtime.GOOS == "darwin" {
			// Look for: open -a "AppName" "$@"
			if strings.Contains(line, "open -a") {
				parts := strings.Split(line, "\"")
				if len(parts) >= 2 {
					return parts[1], nil
				}
			}
		}
	}

	return "unknown", nil
}

func IsInPath() bool {
	dir := GetLauncherDir()
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ":")

	for _, p := range paths {
		if p == dir {
			return true
		}
	}

	return false
}
