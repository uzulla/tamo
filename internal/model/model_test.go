package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewTask(t *testing.T) {
	id := uuid.New().String()
	title := "Test Task"
	description := "Test Description"
	memoRefs := []string{"memo1", "memo2"}

	task := NewTask(id, title, description, memoRefs)

	if task.ID != id {
		t.Errorf("Expected task ID to be %s, got %s", id, task.ID)
	}

	if task.Title != title {
		t.Errorf("Expected task title to be %s, got %s", title, task.Title)
	}

	if task.Description != description {
		t.Errorf("Expected task description to be %s, got %s", description, task.Description)
	}

	if len(task.MemoRefs) != len(memoRefs) {
		t.Errorf("Expected %d memo refs, got %d", len(memoRefs), len(task.MemoRefs))
	}

	for i, ref := range memoRefs {
		if task.MemoRefs[i] != ref {
			t.Errorf("Expected memo ref %d to be %s, got %s", i, ref, task.MemoRefs[i])
		}
	}

	if task.Done {
		t.Error("Expected new task to be not done")
	}

	// Check that CreatedAt and UpdatedAt are set and close to current time
	now := time.Now().UTC()
	if task.CreatedAt.Time.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if task.UpdatedAt.Time.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
	if now.Sub(task.CreatedAt.Time) > 5*time.Second {
		t.Errorf("CreatedAt is too far in the past: %v", task.CreatedAt.Time)
	}
	if now.Sub(task.UpdatedAt.Time) > 5*time.Second {
		t.Errorf("UpdatedAt is too far in the past: %v", task.UpdatedAt.Time)
	}
}

func TestNewMemo(t *testing.T) {
	id := uuid.New().String()
	title := "Test Memo"
	titlePtr := &title
	content := "Test Content"

	memo := NewMemo(id, titlePtr, content)

	if memo.ID != id {
		t.Errorf("Expected memo ID to be %s, got %s", id, memo.ID)
	}

	if *memo.Title != title {
		t.Errorf("Expected memo title to be %s, got %s", title, *memo.Title)
	}

	if memo.Content != content {
		t.Errorf("Expected memo content to be %s, got %s", content, memo.Content)
	}

	// Check that CreatedAt and UpdatedAt are set and close to current time
	now := time.Now().UTC()
	if memo.CreatedAt.Time.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if memo.UpdatedAt.Time.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
	if now.Sub(memo.CreatedAt.Time) > 5*time.Second {
		t.Errorf("CreatedAt is too far in the past: %v", memo.CreatedAt.Time)
	}
	if now.Sub(memo.UpdatedAt.Time) > 5*time.Second {
		t.Errorf("UpdatedAt is too far in the past: %v", memo.UpdatedAt.Time)
	}
}

func TestNewStore(t *testing.T) {
	store := NewStore()

	if store.Version != 1 {
		t.Errorf("Expected store version to be 1, got %d", store.Version)
	}

	if len(store.Tasks) != 0 {
		t.Errorf("Expected empty tasks list, got %d tasks", len(store.Tasks))
	}

	if len(store.Memos) != 0 {
		t.Errorf("Expected empty memos list, got %d memos", len(store.Memos))
	}
}

func TestStore_AddTask(t *testing.T) {
	store := NewStore()
	id := uuid.New().String()
	task := NewTask(id, "Test Task", "Test Description", nil)

	store.AddTask(task)

	if len(store.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(store.Tasks))
	}

	if store.Tasks[0].ID != id {
		t.Errorf("Expected task ID to be %s, got %s", id, store.Tasks[0].ID)
	}
}

func TestStore_AddMemo(t *testing.T) {
	store := NewStore()
	id := uuid.New().String()
	title := "Test Memo"
	titlePtr := &title
	memo := NewMemo(id, titlePtr, "Test Content")

	store.AddMemo(memo)

	if len(store.Memos) != 1 {
		t.Errorf("Expected 1 memo, got %d", len(store.Memos))
	}

	if store.Memos[0].ID != id {
		t.Errorf("Expected memo ID to be %s, got %s", id, store.Memos[0].ID)
	}
}

func TestStore_FindTaskByID(t *testing.T) {
	store := NewStore()
	id := uuid.New().String()
	task := NewTask(id, "Test Task", "Test Description", nil)
	store.AddTask(task)

	// Test finding by full ID
	foundTask := store.FindTaskByID(id)
	if foundTask == nil {
		t.Error("Expected to find task, got nil")
	} else if foundTask.ID != id {
		t.Errorf("Expected task ID to be %s, got %s", id, foundTask.ID)
	}

	// Test not finding a task
	notFoundTask := store.FindTaskByID("nonexistent")
	if notFoundTask != nil {
		t.Errorf("Expected not to find task, got %v", notFoundTask)
	}
}

func TestStore_FindMemoByID(t *testing.T) {
	store := NewStore()
	id := uuid.New().String()
	title := "Test Memo"
	titlePtr := &title
	memo := NewMemo(id, titlePtr, "Test Content")
	store.AddMemo(memo)

	// Test finding by full ID
	foundMemo := store.FindMemoByID(id)
	if foundMemo == nil {
		t.Error("Expected to find memo, got nil")
	} else if foundMemo.ID != id {
		t.Errorf("Expected memo ID to be %s, got %s", id, foundMemo.ID)
	}

	// Test not finding a memo
	notFoundMemo := store.FindMemoByID("nonexistent")
	if notFoundMemo != nil {
		t.Errorf("Expected not to find memo, got %v", notFoundMemo)
	}
}

func TestStore_GetMaxTaskOrder(t *testing.T) {
	store := NewStore()

	// Empty store
	maxOrder := store.GetMaxTaskOrder()
	if maxOrder != 0 {
		t.Errorf("Expected max order to be 0 for empty store, got %f", maxOrder)
	}

	// Add tasks with different orders
	id1 := uuid.New().String()
	task1 := NewTask(id1, "Task 1", "Description 1", nil)
	task1.Order = 1.0
	store.AddTask(task1)

	id2 := uuid.New().String()
	task2 := NewTask(id2, "Task 2", "Description 2", nil)
	task2.Order = 2.5
	store.AddTask(task2)

	id3 := uuid.New().String()
	task3 := NewTask(id3, "Task 3", "Description 3", nil)
	task3.Order = 1.5
	store.AddTask(task3)

	maxOrder = store.GetMaxTaskOrder()
	if maxOrder != 2.5 {
		t.Errorf("Expected max order to be 2.5, got %f", maxOrder)
	}
}

func TestStore_GetMinTaskOrder(t *testing.T) {
	store := NewStore()

	// Empty store
	minOrder := store.GetMinTaskOrder()
	if minOrder != 0 {
		t.Errorf("Expected min order to be 0 for empty store, got %f", minOrder)
	}

	// Add tasks with different orders
	id1 := uuid.New().String()
	task1 := NewTask(id1, "Task 1", "Description 1", nil)
	task1.Order = 1.0
	store.AddTask(task1)

	id2 := uuid.New().String()
	task2 := NewTask(id2, "Task 2", "Description 2", nil)
	task2.Order = 2.5
	store.AddTask(task2)

	id3 := uuid.New().String()
	task3 := NewTask(id3, "Task 3", "Description 3", nil)
	task3.Order = 1.5
	store.AddTask(task3)

	minOrder = store.GetMinTaskOrder()
	if minOrder != 1.0 {
		t.Errorf("Expected min order to be 1.0, got %f", minOrder)
	}
}
