package models

import "time"

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "success"
	StatusFailed     TaskStatus = "failure"
)

type Task struct {
	ID          string
	Command     string
	Payload     string
	TaskStatus  TaskStatus
	CreatedAt   time.Time
	ScheduledAt time.Time
}
