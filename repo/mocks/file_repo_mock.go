package mocks

import (
	"gorm.io/gorm"
	"task-golang/model"

	"github.com/stretchr/testify/mock"
)

type FileRepoMock struct {
	mock.Mock
}

func (r *FileRepoMock) SaveAttachmentFile(attachmentFile model.AttachmentFile) model.AttachmentFile {
	args := r.Called(attachmentFile)
	return args.Get(0).(model.AttachmentFile)
}

func (r *FileRepoMock) SaveTaskImage(taskImage model.TaskImage) model.TaskImage {
	args := r.Called(taskImage)
	return args.Get(0).(model.TaskImage)
}

func (r *FileRepoMock) SaveTaskVideo(taskVideo model.TaskVideo) model.TaskVideo {
	args := r.Called(taskVideo)
	return args.Get(0).(model.TaskVideo)
}

func (r *FileRepoMock) SaveTaskAttachmentFile(taskAttachmentFile *model.TaskAttachmentFile) error {
	args := r.Called(taskAttachmentFile)
	return args.Error(0)
}

func (r *FileRepoMock) SaveTaskTaskImage(taskTaskImage *model.TaskTaskImage) error {
	args := r.Called(taskTaskImage)
	return args.Error(0)
}

func (r *FileRepoMock) SaveTaskTaskVideo(taskTaskVideo *model.TaskTaskVideo) error {
	args := r.Called(taskTaskVideo)
	return args.Error(0)
}

func (r *FileRepoMock) DeleteTaskAttachmentFile(tx *gorm.DB, attachmentFileId int64) error {
	args := r.Called(tx, attachmentFileId)
	return args.Error(0)
}

func (r *FileRepoMock) DeleteAttachmentFile(tx *gorm.DB, attachmentFileId int64) error {
	args := r.Called(tx, attachmentFileId)
	return args.Error(0)
}

func (r *FileRepoMock) FindTaskAttachmentFileByAttachmentFileId(attachmentFileId int64) (*model.TaskAttachmentFile, error) {
	args := r.Called(attachmentFileId)
	if args.Get(0) != nil {
		return args.Get(0).(*model.TaskAttachmentFile), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *FileRepoMock) FindTaskAttachmentsFileByTaskId(taskId int64) (*[]model.TaskAttachmentFile, error) {
	args := r.Called(taskId)
	if args.Get(0) != nil {
		return args.Get(0).(*[]model.TaskAttachmentFile), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *FileRepoMock) FindTaskTaskImageByTaskId(taskId int64) (*model.TaskTaskImage, error) {
	args := r.Called(taskId)
	if args.Get(0) != nil {
		return args.Get(0).(*model.TaskTaskImage), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *FileRepoMock) FindTaskTaskVideo(taskVideoId int64) (*model.TaskTaskVideo, error) {
	args := r.Called(taskVideoId)
	if args.Get(0) != nil {
		return args.Get(0).(*model.TaskTaskVideo), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *FileRepoMock) FindTaskTaskVideosByTaskId(taskId int64) (*[]model.TaskTaskVideo, error) {
	args := r.Called(taskId)
	if args.Get(0) != nil {
		return args.Get(0).(*[]model.TaskTaskVideo), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *FileRepoMock) FindAttachmentFileById(attachmentFileId int64) (*model.AttachmentFile, error) {
	args := r.Called(attachmentFileId)
	if args.Get(0) != nil {
		return args.Get(0).(*model.AttachmentFile), args.Error(1)
	}
	return nil, args.Error(1)
}
