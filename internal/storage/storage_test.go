package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/zishida/tamo/internal/model"
)

func TestStorage_Load(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tamo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a storage with custom paths
	tamoDir := filepath.Join(tempDir, ".tamo")
	dataFile := filepath.Join(tamoDir, "data.json")
	storage := NewStorageWithPath(tamoDir, dataFile)

	// Create the directory
	if err := os.Mkdir(tamoDir, 0755); err != nil {
		t.Fatalf("Failed to create .tamo dir: %v", err)
	}

	// Test loading a non-existent store
	_, err = storage.Load()
	if err == nil {
		t.Error("Expected error when loading non-existent store, got nil")
	}

	// Create an empty store
	emptyStore := model.NewStore()
	if err := storage.Save(emptyStore); err != nil {
		t.Fatalf("Failed to save empty store: %v", err)
	}

	// Test loading the empty store
	loadedStore, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load store: %v", err)
	}

	if loadedStore.Version != emptyStore.Version {
		t.Errorf("Expected version %d, got %d", emptyStore.Version, loadedStore.Version)
	}

	if len(loadedStore.Tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(loadedStore.Tasks))
	}

	if len(loadedStore.Memos) != 0 {
		t.Errorf("Expected 0 memos, got %d", len(loadedStore.Memos))
	}

	// Create a store with tasks and memos
	store := model.NewStore()

	taskID := uuid.New().String()
	task := model.NewTask(taskID, "Test Task", "Test Description", nil)
	task.Order = 1.0
	store.AddTask(task)

	memoID := uuid.New().String()
	title := "Test Memo"
	memo := model.NewMemo(memoID, &title, "Test Content")
	store.AddMemo(memo)

	// Save the store
	if err := storage.Save(store); err != nil {
		t.Fatalf("Failed to save store: %v", err)
	}

	// Load the store
	loadedStore, err = storage.Load()
	if err != nil {
		t.Fatalf("Failed to load store: %v", err)
	}

	// Verify the loaded store
	if len(loadedStore.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(loadedStore.Tasks))
	} else {
		loadedTask := loadedStore.Tasks[0]
		if loadedTask.ID != taskID {
			t.Errorf("Expected task ID %s, got %s", taskID, loadedTask.ID)
		}
		if loadedTask.Title != "Test Task" {
			t.Errorf("Expected task title 'Test Task', got '%s'", loadedTask.Title)
		}
		if loadedTask.Description != "Test Description" {
			t.Errorf("Expected task description 'Test Description', got '%s'", loadedTask.Description)
		}
		if loadedTask.Order != 1.0 {
			t.Errorf("Expected task order 1.0, got %f", loadedTask.Order)
		}
	}

	if len(loadedStore.Memos) != 1 {
		t.Errorf("Expected 1 memo, got %d", len(loadedStore.Memos))
	} else {
		loadedMemo := loadedStore.Memos[0]
		if loadedMemo.ID != memoID {
			t.Errorf("Expected memo ID %s, got %s", memoID, loadedMemo.ID)
		}
		if *loadedMemo.Title != "Test Memo" {
			t.Errorf("Expected memo title 'Test Memo', got '%s'", *loadedMemo.Title)
		}
		if loadedMemo.Content != "Test Content" {
			t.Errorf("Expected memo content 'Test Content', got '%s'", loadedMemo.Content)
		}
	}
}

func TestStorage_Initialize(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tamo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a storage with custom paths
	tamoDir := filepath.Join(tempDir, ".tamo")
	dataFile := filepath.Join(tamoDir, "data.json")
	storage := NewStorageWithPath(tamoDir, dataFile)

	// Test initializing a store
	if err := storage.Initialize(); err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	// Check that the .tamo directory was created
	if _, err := os.Stat(tamoDir); os.IsNotExist(err) {
		t.Errorf("Expected .tamo directory to exist, but it doesn't")
	}

	// Check that the data.json file was created
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		t.Errorf("Expected data.json file to exist, but it doesn't")
	}

	// Test initializing a store that already exists (should not error)
	if err := storage.Initialize(); err != nil {
		t.Errorf("Unexpected error when initializing existing store: %v", err)
	}
}

func TestStorage_Save(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tamo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a storage with custom paths
	tamoDir := filepath.Join(tempDir, ".tamo")
	dataFile := filepath.Join(tamoDir, "data.json")
	storage := NewStorageWithPath(tamoDir, dataFile)

	// Create the directory
	if err := os.Mkdir(tamoDir, 0755); err != nil {
		t.Fatalf("Failed to create .tamo dir: %v", err)
	}

	// Create a store with tasks and memos
	store := model.NewStore()

	taskID := uuid.New().String()
	task := model.NewTask(taskID, "Test Task", "Test Description", nil)
	task.Order = 1.0
	store.AddTask(task)

	memoID := uuid.New().String()
	title := "Test Memo"
	memo := model.NewMemo(memoID, &title, "Test Content")
	store.AddMemo(memo)

	// Save the store
	if err := storage.Save(store); err != nil {
		t.Fatalf("Failed to save store: %v", err)
	}

	// Check that the data.json file was created
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		t.Errorf("Expected data.json file to exist, but it doesn't")
	}

	// Load the store and verify it
	loadedStore, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load store: %v", err)
	}

	if len(loadedStore.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(loadedStore.Tasks))
	}

	if len(loadedStore.Memos) != 1 {
		t.Errorf("Expected 1 memo, got %d", len(loadedStore.Memos))
	}
}

func TestStorage_Exists(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tamo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a storage with custom paths
	tamoDir := filepath.Join(tempDir, ".tamo")
	dataFile := filepath.Join(tamoDir, "data.json")
	storage := NewStorageWithPath(tamoDir, dataFile)

	// Test exists on non-existent file
	if storage.Exists() {
		t.Error("Expected Exists() to return false for non-existent file, got true")
	}

	// Create the directory and file
	if err := os.Mkdir(tamoDir, 0755); err != nil {
		t.Fatalf("Failed to create .tamo dir: %v", err)
	}
	if err := os.WriteFile(dataFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create data.json file: %v", err)
	}

	// Test exists on existing file
	if !storage.Exists() {
		t.Error("Expected Exists() to return true for existing file, got false")
	}
}

func TestStorage_EnsureDirectoryExists(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tamo-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a storage with custom paths
	tamoDir := filepath.Join(tempDir, ".tamo")
	dataFile := filepath.Join(tamoDir, "data.json")
	storage := NewStorageWithPath(tamoDir, dataFile)

	// Test ensuring directory exists when it doesn't
	if err := storage.EnsureDirectoryExists(); err != nil {
		t.Fatalf("Failed to ensure directory exists: %v", err)
	}

	// Check that the directory was created
	if _, err := os.Stat(tamoDir); os.IsNotExist(err) {
		t.Errorf("Expected directory to exist, but it doesn't")
	}

	// Test ensuring directory exists when it already does
	if err := storage.EnsureDirectoryExists(); err != nil {
		t.Fatalf("Failed to ensure directory exists when it already does: %v", err)
	}
}
