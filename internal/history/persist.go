package history

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultFileName = "cronlens_history.json"

// SaveToFile writes the current history entries to a JSON file.
func (h *History) SaveToFile(path string) error {
	if path == "" {
		path = defaultFilePath()
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	h.mu.Lock()
	data, err := json.MarshalIndent(h.entries, "", "  ")
	h.mu.Unlock()
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// LoadFromFile reads history entries from a JSON file into the current History.
// If the file does not exist, it returns nil without error.
func (h *History) LoadFromFile(path string) error {
	if path == "" {
		path = defaultFilePath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	var entries []*Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, e := range entries {
		if len(h.entries) >= h.maxSize {
			break
		}
		h.entries = append(h.entries, e)
	}
	return nil
}

func defaultFilePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "cronlens", defaultFileName)
}
