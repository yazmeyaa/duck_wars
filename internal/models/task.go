package models

import "time"

type TaskStatus string

const (
	STATUS_PENDING    TaskStatus = "pending"
	STATUS_PROCESSING TaskStatus = "processing"
	STATUS_COMPLETED  TaskStatus = "success"
	STATUS_FAILDED    TaskStatus = "failure"
)

type Task struct {
	ID         string
	Command    string
	Payload    string
	TaskStatus TaskStatus
	CreatedAt  time.Time
	SheduledAt time.Time
}
