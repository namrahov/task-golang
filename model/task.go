package model

import "time"

type Task struct {
	Id       int64    `gorm:"primaryKey;column:id" json:"id"`
	Name     string   `gorm:"column:name" json:"name"`
	Priority Priority `gorm:"column:priority" json:"priority"`
	Status   Status   `gorm:"column:status" json:"status"`

	// The ID columns:
	CreatedByID  *int64 `gorm:"column:created_by_id"`  // new field
	ChangedByID  *int64 `gorm:"column:changed_by_id"`  // new field
	AssignedByID *int64 `gorm:"column:assigned_by_id"` // new field
	AssignedToID *int64 `gorm:"column:assigned_to_id"` // new field
	BoardID      *int64 `gorm:"column:board_id"`       // new field

	// The actual relationship fields:
	CreatedBy  *User  `gorm:"foreignKey:CreatedByID" json:"createdBy"`
	ChangedBy  *User  `gorm:"foreignKey:ChangedByID" json:"changedBy"`
	AssignedBy *User  `gorm:"foreignKey:AssignedByID" json:"assignedBy"`
	AssignedTo *User  `gorm:"foreignKey:AssignedToID" json:"assignedTo"`
	Board      *Board `gorm:"foreignKey:BoardID"       json:"board"`

	Deadline  time.Time `gorm:"column:deadline" json:"deadline"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
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

type TaskRequestDto struct {
	Name     string
	Priority Priority
	Deadline string
}
