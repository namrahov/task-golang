package repo

import (
	"gorm.io/gorm"
	"task-golang/model"
)

type IFileRepo interface {
	SaveAttachmentFile(attachmentFile model.AttachmentFile) model.AttachmentFile
	SaveTaskImage(taskImage model.TaskImage) model.TaskImage
	SaveTaskVideo(taskVideo model.TaskVideo) model.TaskVideo
	SaveTaskAttachmentFile(taskAttachmentFile *model.TaskAttachmentFile) error
	SaveTaskTaskImage(taskTaskImage *model.TaskTaskImage) error
	SaveTaskTaskVideo(taskTaskVideo *model.TaskTaskVideo) error
	DeleteTaskAttachmentFile(tx *gorm.DB, attachmentFileId int64) error
	FindTaskAttachmentFileByAttachmentFileId(attachmentFileId int64) (*model.TaskAttachmentFile, error)
	FindTaskTaskImageByTaskId(taskId int64) (*model.TaskTaskImage, error)
	FindTaskTaskVideoByTaskId(taskId int64) (*model.TaskTaskVideo, error)
	DeleteAttachmentFile(tx *gorm.DB, attachmentFileId int64) error
	FindAttachmentFileById(attachmentFileId int64) (*model.AttachmentFile, error)
}

type FileRepo struct {
}

func (fr *FileRepo) SaveAttachmentFile(attachmentFile model.AttachmentFile) model.AttachmentFile {
	// Save the file record to the database
	result := Db.Create(&attachmentFile)
	if result.Error != nil {
		// Handle the error (you can log it or return an empty file object)
		panic(result.Error)
	}

	return attachmentFile
}

func (fr *FileRepo) SaveTaskImage(taskImage model.TaskImage) model.TaskImage {
	result := Db.Create(&taskImage)
	if result.Error != nil {
		// Handle the error (you can log it or return an empty file object)
		panic(result.Error)
	}

	return taskImage
}

func (fr *FileRepo) SaveTaskVideo(taskVideo model.TaskVideo) model.TaskVideo {
	result := Db.Create(&taskVideo)
	if result.Error != nil {
		panic(result.Error)
	}

	return taskVideo
}

func (fr *FileRepo) SaveTaskAttachmentFile(taskAttachmentFile *model.TaskAttachmentFile) error {
	result := Db.Create(taskAttachmentFile)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (fr *FileRepo) SaveTaskTaskImage(taskTaskImage *model.TaskTaskImage) error {
	result := Db.Create(taskTaskImage)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (fr *FileRepo) SaveTaskTaskVideo(taskTaskVideo *model.TaskTaskVideo) error {
	result := Db.Create(taskTaskVideo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (fr *FileRepo) DeleteTaskAttachmentFile(tx *gorm.DB, attachmentFileId int64) error {
	result := tx.Where("attachment_file_id = ?", attachmentFileId).Delete(&model.TaskAttachmentFile{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (fr *FileRepo) DeleteAttachmentFile(tx *gorm.DB, attachmentFileId int64) error {
	result := tx.Where("id = ?", attachmentFileId).Delete(&model.AttachmentFile{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (fr *FileRepo) FindTaskAttachmentFileByAttachmentFileId(attachmentFileId int64) (*model.TaskAttachmentFile, error) {
	var taskAttachmentFile model.TaskAttachmentFile

	result := Db.Preload("AttachmentFile").Where("attachment_file_id = ?", attachmentFileId).First(&taskAttachmentFile)
	if result.Error != nil {
		return nil, result.Error
	}

	return &taskAttachmentFile, nil
}

func (fr *FileRepo) FindTaskTaskImageByTaskId(taskId int64) (*model.TaskTaskImage, error) {
	var taskTaskImage model.TaskTaskImage

	result := Db.Preload("TaskImage").Where("task_id = ?", taskId).First(&taskTaskImage)
	if result.Error != nil {
		return nil, result.Error
	}

	return &taskTaskImage, nil
}

func (fr *FileRepo) FindTaskTaskVideoByTaskId(taskId int64) (*model.TaskTaskVideo, error) {
	var taskTaskVideo model.TaskTaskVideo

	result := Db.Preload("TaskVideo").Where("task_id = ?", taskId).First(&taskTaskVideo)
	if result.Error != nil {
		return nil, result.Error
	}

	return &taskTaskVideo, nil
}

func (fr *FileRepo) FindAttachmentFileById(attachmentFileId int64) (*model.AttachmentFile, error) {
	var attachmentFile model.AttachmentFile // Use the correct type

	// Query the database
	result := Db.Where("id = ?", attachmentFileId).First(&attachmentFile)
	if result.Error != nil {
		return nil, result.Error // Return nil and the error if the query fails
	}

	return &attachmentFile, nil // Return a pointer to the found record
}
