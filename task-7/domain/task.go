package domain

import (
	"errors"
	"time"
)

type TaskStatus string

const (
	StatusPending TaskStatus = "pending"
	StatusDoing   TaskStatus = "doing"
	StatusDone    TaskStatus = "done"
)

var validTaskStatuses = map[TaskStatus]bool{
	StatusPending: true,
	StatusDoing:   true,
	StatusDone:    true,
}

func IsTaskStatusValid(ts string) bool {
	_, ok := validTaskStatuses[TaskStatus(ts)]
	return ok
}

func NewTaskStatusFromString(s string) (TaskStatus, error) {
	status := s
	if !IsTaskStatusValid(status) {
		return "", errors.New("invalid task status")
	}
	return TaskStatus(status), nil
}

type Task struct {
	TaskID      string
	Title       string
	Description string
	Status      TaskStatus
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Task) MarkCompleted() {
	t.Status = StatusDone
	t.UpdatedAt = time.Now().UTC()
}

func (t *Task) IsOverDue() bool {
	return time.Now().UTC().After(t.DueDate) && t.Status != StatusDone
}
