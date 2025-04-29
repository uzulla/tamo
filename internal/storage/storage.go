package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/zishida/tamo/internal/model"
)

const (
	// DefaultDirName is the default directory name for tamo data
	DefaultDirName = ".tamo"
	// DefaultFileName is the default file name for tamo data
	DefaultFileName = "data.json"
)

// Storage handles the persistence of the store
type Storage struct {
	DirPath  string
	FilePath string
}

// NewStorage creates a new storage with the default path
func NewStorage() *Storage {
	return &Storage{
		DirPath:  DefaultDirName,
		FilePath: filepath.Join(DefaultDirName, DefaultFileName),
	}
}

// NewStorageWithPath creates a new storage with the given path
func NewStorageWithPath(dirPath, filePath string) *Storage {
	return &Storage{
		DirPath:  dirPath,
		FilePath: filePath,
	}
}

// Initialize creates the directory and empty data file if they don't exist
func (s *Storage) Initialize() error {
	// Check if directory exists
	if _, err := os.Stat(s.DirPath); os.IsNotExist(err) {
		// Create directory
		if err := os.Mkdir(s.DirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Check if file exists
	if _, err := os.Stat(s.FilePath); os.IsNotExist(err) {
		// Create empty store
		store := model.NewStore()

		// Save empty store
		if err := s.Save(store); err != nil {
			return fmt.Errorf("failed to create empty data file: %w", err)
		}
	}

	return nil
}

// Load loads the store from the file
func (s *Storage) Load() (*model.Store, error) {
	// Check if file exists
	if _, err := os.Stat(s.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("data file not found: %s", s.FilePath)
	}

	// Read file
	data, err := ioutil.ReadFile(s.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	// Parse JSON
	var store model.Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to parse data file: %w", err)
	}

	// Fix time fields
	for _, task := range store.Tasks {
		if task.CreatedAt.IsZero() {
			task.CreatedAt = model.CustomTime{time.Now().UTC()}
		}
		if task.UpdatedAt.IsZero() {
			task.UpdatedAt = model.CustomTime{time.Now().UTC()}
		}
	}
	for _, memo := range store.Memos {
		if memo.CreatedAt.IsZero() {
			memo.CreatedAt = model.CustomTime{time.Now().UTC()}
		}
		if memo.UpdatedAt.IsZero() {
			memo.UpdatedAt = model.CustomTime{time.Now().UTC()}
		}
	}

	return &store, nil
}

// Save saves the store to the file atomically
func (s *Storage) Save(store *model.Store) error {
	// Marshal JSON
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Create temporary file
	tmpFile, err := ioutil.TempFile(s.DirPath, "data.*.json.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write data to temporary file
	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}

	// Close temporary file
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Rename temporary file to target file (atomic operation)
	if err := os.Rename(tmpFile.Name(), s.FilePath); err != nil {
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// Exists checks if the data file exists
func (s *Storage) Exists() bool {
	_, err := os.Stat(s.FilePath)
	return err == nil
}

// EnsureDirectoryExists ensures that the directory exists
func (s *Storage) EnsureDirectoryExists() error {
	if _, err := os.Stat(s.DirPath); os.IsNotExist(err) {
		return os.Mkdir(s.DirPath, 0755)
	}
	return nil
}
