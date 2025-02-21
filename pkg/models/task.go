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
	ID          string     `json:"id"`
	Command     string     `json:"command"`
	Payload     string     `json:"payload"`
	TaskStatus  TaskStatus `json:"taskStatus"`
	CreatedAt   time.Time  `json:"createdAt"`
	ScheduledAt time.Time  `json:"scheduledAt"`
}
