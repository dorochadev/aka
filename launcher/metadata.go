package launcher

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const metadataFile = ".config/aka/launchers.json"

type MetadataStore map[string]*LauncherMetadata

func getMetadataPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, metadataFile), nil
}

func LoadMetadata() (MetadataStore, error) {
	path, err := getMetadataPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return make(MetadataStore), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var store MetadataStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return store, nil
}

func SaveMetadata(store MetadataStore) error {
	path, err := getMetadataPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

func GetMetadata(name string) (*LauncherMetadata, error) {
	store, err := LoadMetadata()
	if err != nil {
		return nil, err
	}
	return store[name], nil
}

func SetMetadata(name string, metadata *LauncherMetadata) error {
	store, err := LoadMetadata()
	if err != nil {
		return err
	}

	store[name] = metadata
	return SaveMetadata(store)
}

func DeleteMetadata(name string) error {
	store, err := LoadMetadata()
	if err != nil {
		return err
	}

	delete(store, name)
	return SaveMetadata(store)
}
