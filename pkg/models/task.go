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
	ID         string     `json:"id"`
	Command    string     `json:"command"`
	Payload    string     `json:"payload"`
	TaskStatus TaskStatus `json:"taskStatus"`
	CreatedAt  time.Time  `json:"createdAt"`
	SheduledAt time.Time  `json:"sheduledAt"`
}
