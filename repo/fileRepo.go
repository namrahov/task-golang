package repo

import (
	"task-golang/model"
)

type IFileRepo interface {
	SaveAttachmentFile(attachmentFile model.AttachmentFile) model.AttachmentFile
	SaveTaskAttachmentFile(taskAttachmentFile *model.TaskAttachmentFile) error
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
