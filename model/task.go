package model

import "time"

type Task struct {
	Id       int64    `gorm:"primaryKey;column:id" json:"id"`
	Name     string   `gorm:"column:name" json:"name"`
	Priority Priority `gorm:"column:priority" json:"priority"`
	Status   Status   `gorm:"column:status" json:"status"`

	// The ID columns:
	CreatedByID  *int64 `gorm:"column:created_by_id"`
	ChangedByID  *int64 `gorm:"column:changed_by_id"`
	AssignedByID *int64 `gorm:"column:assigned_by_id"`
	AssignedToID *int64 `gorm:"column:assigned_to_id"`
	BoardID      *int64 `gorm:"column:board_id"`

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
	Name     string   `json:"name"`
	Priority Priority `json:"priority"`
	Deadline string   `json:"deadline"`
}

type TaskResponseDto struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	Priority Priority `json:"priority"`
	Status   Status   `json:"status"`

	CreatedBy  *User  `json:"createdBy"`
	ChangedBy  *User  `json:"changedBy"`
	AssignedBy *User  `json:"assignedBy"`
	AssignedTo *User  `json:"assignedTo"`
	Board      *Board `json:"board"`

	Deadline time.Time `json:"deadline"`

	AttachmentFileIds []int64 `json:"attachmentFileId"`
	TaskVideoId       []int64 `json:"taskVideoId"`
	TaskImageUrl      string  `json:"taskImageUrl"`
}
