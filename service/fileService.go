package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"strconv"
	"task-golang/config"
	"task-golang/mapper"
	"task-golang/model"
	"task-golang/repo"
	"task-golang/util"
)

type IFileService interface {
	UploadAttachmentFile(ctx context.Context, multipartFile *multipart.File, multipartFileHeader *multipart.FileHeader, taskId int64) (*model.FileResponseDto, *model.ErrorResponse)
}

type FileService struct {
	FileRepo repo.IFileRepo
	UserRepo repo.IUserRepo
	UserUtil util.IUserUtil
}

func (fs *FileService) UploadAttachmentFile(ctx context.Context, multipartFile *multipart.File, multipartFileHeader *multipart.FileHeader, taskId int64) (*model.FileResponseDto, *model.ErrorResponse) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.UploadAttachmentFile.start")

	// Check if competitionId or multipartFile is nil
	if taskId <= 0 || multipartFileHeader == nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-continue", model.Exception),
			Message: "Task ID and file must not be null or empty",
			Code:    http.StatusBadRequest,
		}
	}

	// Check if file size exceeds limit
	maxSize, err := strconv.ParseInt(config.Props.AttachmentFileMaxSize, 10, 64)
	if err != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-continue", model.Exception),
			Message: "Invalid max file size configuration",
			Code:    http.StatusBadRequest,
		}
	}

	if multipartFileHeader.Size > maxSize {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-continue", model.Exception),
			Message: "Image is too large. Maximum size is " + config.Props.AttachmentFileMaxSize + " bytes",
			Code:    http.StatusBadRequest,
		}
	}

	user, errGetUser := fs.UserUtil.GetUserFromRequest(ctx)
	if errGetUser != nil {
		return nil, errGetUser
	}

	attachmentFileDto, errBuildAttachmentFile := mapper.BuildAttachmentFile(multipartFileHeader, config.Props.MinioBucket)
	if errBuildAttachmentFile != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-build-file", model.Exception),
			Message: errBuildAttachmentFile.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	attachmentFile := fs.FileRepo.SaveAttachmentFile(attachmentFileDto.AttachmentFile)

	taskAttachmentFile := &model.TaskAttachmentFile{
		UserID:           &user.Id,
		TaskID:           &taskId,
		AttachmentFileID: &attachmentFile.Id,
	}

	errSaveTaskAttachmentFile := fs.FileRepo.SaveTaskAttachmentFile(taskAttachmentFile)
	if errSaveTaskAttachmentFile != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-save-taskAttachmentFile", model.Exception),
			Message: errBuildAttachmentFile.Error(),
			Code:    http.StatusForbidden,
		}
	}

	// Initialize Minio client
	minioClient, err := config.NewMinioClient()
	if err != nil {
		log.Fatalf("Failed to initialize Minio client: %v", err)
	}

	errUploadFileToMinio := util.UploadFileToMinio(ctx, attachmentFileDto.UniqueName, *multipartFile, multipartFileHeader.Size, minioClient)
	if errUploadFileToMinio != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-uplaod-to-minio", model.Exception),
			Message: errUploadFileToMinio.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	logger.Info("ActionLog.UploadAttachmentFile.end")
	return &model.FileResponseDto{attachmentFile.Id}, nil
}
