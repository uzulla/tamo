package model

import (
	"encoding/json"
	"time"
)

// CustomTime is a wrapper around time.Time that formats as ISO 8601 in JSON
type CustomTime struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface
func (t CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(time.RFC3339))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *CustomTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return err
	}

	t.Time = parsedTime
	return nil
}

// Task represents a task to be done with properties like ID, title, description, order, completion status, and memo references
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Order       float64    `json:"order"`
	Done        bool       `json:"done"`
	MemoRefs    []string   `json:"memo_refs"`
	CreatedAt   CustomTime `json:"created_at"`
	UpdatedAt   CustomTime `json:"updated_at"`
}

// Memo stores information related to tasks with properties like ID, title, and content
type Memo struct {
	ID        string     `json:"id"`
	Title     *string    `json:"title"` // Optional
	Content   string     `json:"content"`
	CreatedAt CustomTime `json:"created_at"`
	UpdatedAt CustomTime `json:"updated_at"`
}

// Store is the main data structure that contains all tasks and memos
type Store struct {
	Version int     `json:"version"`
	Tasks   []*Task `json:"tasks"`
	Memos   []*Memo `json:"memos"`
}

// NewStore creates a new empty store with version 1
func NewStore() *Store {
	return &Store{
		Version: 1,
		Tasks:   make([]*Task, 0),
		Memos:   make([]*Memo, 0),
	}
}

// NewTask creates a new task with the given title, description, and memo references
func NewTask(id, title, description string, memoRefs []string) *Task {
	now := CustomTime{time.Now().UTC()}
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Order:       0.0, // Will be set by the caller
		Done:        false,
		MemoRefs:    memoRefs,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewMemo creates a new memo with the given title and content
func NewMemo(id string, title *string, content string) *Memo {
	now := CustomTime{time.Now().UTC()}
	return &Memo{
		ID:        id,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GetMaxTaskOrder returns the maximum order value of all tasks in the store
func (s *Store) GetMaxTaskOrder() float64 {
	maxOrder := 0.0
	for _, task := range s.Tasks {
		if task.Order > maxOrder {
			maxOrder = task.Order
		}
	}
	return maxOrder
}

// GetMinTaskOrder returns the minimum order value of all tasks in the store
func (s *Store) GetMinTaskOrder() float64 {
	if len(s.Tasks) == 0 {
		return 0.0
	}

	minOrder := s.Tasks[0].Order
	for _, task := range s.Tasks {
		if task.Order < minOrder {
			minOrder = task.Order
		}
	}
	return minOrder
}

// FindTaskByID returns a task by its ID
func (s *Store) FindTaskByID(id string) *Task {
	for _, task := range s.Tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

// FindMemoByID returns a memo by its ID
func (s *Store) FindMemoByID(id string) *Memo {
	for _, memo := range s.Memos {
		if memo.ID == id {
			return memo
		}
	}
	return nil
}

// AddTask adds a task to the store
func (s *Store) AddTask(task *Task) {
	s.Tasks = append(s.Tasks, task)
}

// AddMemo adds a memo to the store
func (s *Store) AddMemo(memo *Memo) {
	s.Memos = append(s.Memos, memo)
}
