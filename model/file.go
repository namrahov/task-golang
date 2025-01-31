package model

import "time"

type TaskAttachmentFile struct {
	Id int64 `gorm:"primaryKey;column:id" json:"id"`
	// The ID columns:
	UserID           *int64 `gorm:"column:user_id"`
	TaskID           *int64 `gorm:"column:task_id"`
	AttachmentFileID *int64 `gorm:"column:attachment_file_id"`

	// The actual relationship fields:
	CreatedBy      *User           `gorm:"foreignKey:UserID" json:"createdBy"`
	Task           *Task           `gorm:"foreignKey:TaskID" json:"task"`
	AttachmentFile *AttachmentFile `gorm:"foreignKey:AttachmentFileID" json:"attachmentFile"`
}

// TableName overrides the default table name
func (TaskAttachmentFile) TableName() string {
	return "task_attachment_file"
}

type AttachmentFile struct {
	Id        int64     `gorm:"primaryKey;column:id" json:"id"`
	FileType  string    `gorm:"column:file_type" json:"fileType"`
	FilePath  string    `gorm:"column:file_path" json:"filePath"`
	FileName  string    `gorm:"column:file_name" json:"fileName"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
}

// TableName overrides the default table name
func (AttachmentFile) TableName() string {
	return "attachment_file"
}

type FileResponseDto struct {
	Id int64 `json:"id"`
}

type TaskAttachmentFileDto struct {
	attachmentFileId int64
	fileName         string
	extension        string
}

type AttachmentFileDto struct {
	AttachmentFile AttachmentFile `json:"attachmentFile"`
	UniqueName     string         `json:"uniqueName"`
}

type TaskImage struct {
	Id        int64     `gorm:"primaryKey;column:id" json:"id"`
	FileType  string    `gorm:"column:file_type" json:"fileType"`
	FilePath  string    `gorm:"column:file_path" json:"filePath"`
	FileName  string    `gorm:"column:file_name" json:"fileName"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
}

// TableName overrides the default table name
func (TaskImage) TableName() string {
	return "task_image"
}

type TaskTaskImage struct {
	TaskID      int64 `gorm:"column:task_id;primaryKey;not null"`
	TaskImageID int64 `gorm:"column:task_image_id;primaryKey;not null"`

	Task      *Task      `gorm:"foreignKey:TaskID" json:"task"`
	TaskImage *TaskImage `gorm:"foreignKey:TaskImageID" json:"taskImage"`
}

// TableName overrides the default table name for GORM
func (TaskTaskImage) TableName() string {
	return "task_task_image"
}

type TaskImageDto struct {
	TaskImage  TaskImage `json:"taskImage"`
	UniqueName string    `json:"uniqueName"`
}

type TaskVideo struct {
	Id        int64     `gorm:"primaryKey;column:id" json:"id"`
	FileType  string    `gorm:"column:file_type" json:"fileType"`
	FilePath  string    `gorm:"column:file_path" json:"filePath"`
	FileName  string    `gorm:"column:file_name" json:"fileName"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
}

// TableName overrides the default table name
func (TaskVideo) TableName() string {
	return "task_video"
}

type TaskTaskVideo struct {
	TaskID      int64 `gorm:"column:task_id;primaryKey;not null"`
	TaskVideoID int64 `gorm:"column:task_video_id;primaryKey;not null"`

	Task      *Task      `gorm:"foreignKey:TaskID" json:"task"`
	TaskVideo *TaskVideo `gorm:"foreignKey:TaskVideoID" json:"taskVideo"`
}

// TableName overrides the default table name for GORM
func (TaskTaskVideo) TableName() string {
	return "task_task_video"
}

type TaskVideoDto struct {
	TaskVideo  TaskVideo `json:"taskVideo"`
	UniqueName string    `json:"uniqueName"`
}
