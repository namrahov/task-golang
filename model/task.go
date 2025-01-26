package model

import "time"

type Task struct {
	Id         int64     `gorm:"primaryKey;column:id" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	Priority   Priority  `gorm:"column:priority" json:"priority"`
	Status     Status    `gorm:"column:status" json:"status"`
	CreatedBy  *User     `gorm:"foreignKey:UserId;references:Id" json:"createdBy"` // Assuming this references the creator's ID
	ChangedBy  *User     `gorm:"foreignKey:UserId;references:Id" json:"changedBy"` // Assuming this references the changer's ID
	AssignedBy *User     `gorm:"foreignKey:UserId;references:Id" json:"assignedBy"`
	AssignedTo *User     `gorm:"foreignKey:UserId;references:Id" json:"assignedTo"` // Many-to-One relationship
	Board      *Board    `gorm:"foreignKey:BoardId;references:Id" json:"board"`     // Many-to-One relationship
	Deadline   time.Time `gorm:"column:deadline" json:"deadline"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"-"`
}

type Status string

const (
	NotStarted   Status = "NOT_STARTED"
	InProgress   Status = "IN_PROGRESS"
	ReadyForTest Status = "READY_FOR_TEST"
	Done         Status = "DONE"
)

type Priority string

const (
	High   Priority = "HIGH"
	Medium Priority = "MEDIUM"
	Low    Priority = "LOW"
)
