package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Status      string             `json:"status" bson:"status"`
	DueDate     *time.Time         `json:"due_date,omitempty" bson:"due_date,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
)

func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("title is required")
	}
	if t.Status == "" {
		t.Status = StatusPending
	}
	if t.Status != StatusPending && t.Status != StatusInProgress && t.Status != StatusCompleted {
		return errors.New("status must be pending, in_progress, or completed")
	}
	return nil
}
