package repo

import (
	"gorm.io/gorm"
	"task-golang/model"
)

type IFileRepo interface {
	SaveAttachmentFile(attachmentFile model.AttachmentFile) model.AttachmentFile
	SaveTaskAttachmentFile(taskAttachmentFile *model.TaskAttachmentFile) error
	DeleteTaskAttachmentFile(tx *gorm.DB, attachmentFileId int64) error
	FindTaskAttachmentFileByAttachmentFileId(attachmentFileId int64) (*model.TaskAttachmentFile, error)
	DeleteAttachmentFile(tx *gorm.DB, attachmentFileId int64) error
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

func (fr *FileRepo) SaveTaskAttachmentFile(taskAttachmentFile *model.TaskAttachmentFile) error {
	result := Db.Create(taskAttachmentFile)
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
