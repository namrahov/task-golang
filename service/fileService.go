package service

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"task-golang/config"
	"task-golang/mapper"
	"task-golang/model"
	"task-golang/repo"
	"task-golang/util"
)

type IFileService interface {
	UploadAttachmentFile(ctx context.Context, multipartFile *multipart.File, multipartFileHeader *multipart.FileHeader, taskId int64) (*model.FileResponseDto, *model.ErrorResponse)
	DeleteAttachmentFile(ctx context.Context, attachmentFileId int64) *model.ErrorResponse
	DownloadAttachmentFile(ctx context.Context, attachmentFileId int64, w http.ResponseWriter) *model.ErrorResponse
	UploadTaskImage(ctx context.Context, multipartFile *multipart.File, multipartFileHeader *multipart.FileHeader, taskId int64) (*model.FileResponseDto, *model.ErrorResponse)
	UploadTaskVideo(ctx context.Context, multipartFile *multipart.File, multipartFileHeader *multipart.FileHeader, taskId int64) (*model.FileResponseDto, *model.ErrorResponse)
	GetTaskImage(ctx context.Context, taskId int64, w http.ResponseWriter) *model.ErrorResponse
	StreamTaskVideo(r *http.Request, taskId int64, w http.ResponseWriter) *model.ErrorResponse
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

func (fs *FileService) DeleteAttachmentFile(ctx context.Context, attachmentFileId int64) *model.ErrorResponse {
	// Begin GORM transaction
	tx := repo.BeginTransaction()
	if tx.Error != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.transaction-begin-failed", model.Exception),
			Message: tx.Error.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	taskAttachmentFile, errFind := fs.FileRepo.FindTaskAttachmentFileByAttachmentFileId(attachmentFileId)
	if errFind != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-find-task-attachment-file", model.Exception),
			Message: errFind.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	errDeleteTask := fs.FileRepo.DeleteTaskAttachmentFile(tx, attachmentFileId)
	if errDeleteTask != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-delete-task-attachment-file", model.Exception),
			Message: errDeleteTask.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	errDeleteAttachmentFile := fs.FileRepo.DeleteAttachmentFile(tx, attachmentFileId)
	if errDeleteAttachmentFile != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-delete-attachment-file", model.Exception),
			Message: errDeleteAttachmentFile.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	minioClient, errMinioClient := config.NewMinioClient()
	if errMinioClient != nil {
		log.Fatalf("Failed to initialize Minio client: %v", errMinioClient)
	}

	errDelete := util.DeleteFileFromMinio(ctx, taskAttachmentFile.AttachmentFile.FilePath, minioClient)
	if errDelete != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-delete-file-from-minio", model.Exception),
			Message: errDelete.Error(),
			Code:    http.StatusForbidden,
		}
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else {
			if errDeleteTask != nil || errDeleteAttachmentFile != nil || errDelete != nil || errMinioClient != nil {
				_ = tx.Rollback()
			} else {
				if err := tx.Commit().Error; err != nil {
					fmt.Println("Transaction commit failed:", err)
				}
			}
		}
	}()

	return nil
}

func (fs *FileService) DownloadAttachmentFile(ctx context.Context, attachmentFileId int64, w http.ResponseWriter) *model.ErrorResponse {
	attachmentFile, errFindAttachmentFile := fs.FileRepo.FindAttachmentFileById(attachmentFileId)
	if errFindAttachmentFile != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-find-attachement-file", model.Exception),
			Message: errFindAttachmentFile.Error(),
			Code:    http.StatusForbidden,
		}
	}

	filePath := attachmentFile.FilePath
	// Extract the object name from the file path
	objectName := filePath[strings.LastIndex(filePath, "/")+1:]

	minioClient, errMinioClient := config.NewMinioClient()
	if errMinioClient != nil {
		log.Fatalf("Failed to initialize Minio client: %v", errMinioClient)
	}
	// Download the object from MinIO
	object, err := minioClient.GetObject(context.Background(), config.Props.MinioBucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-get-file-from-minio", model.Exception),
			Message: err.Error(),
			Code:    http.StatusForbidden,
		}
	}
	defer func(object *minio.Object) {
		err := object.Close()
		if err != nil {
			fmt.Println("cant get object")
		}
	}(object)

	// Set the response headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName))

	// Stream the file content to the response
	_, err = io.Copy(w, object)
	if err != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-download", model.Exception),
			Message: err.Error(),
			Code:    http.StatusForbidden,
		}
	}

	return nil
}

func (fs *FileService) UploadTaskImage(ctx context.Context, multipartFile *multipart.File, multipartFileHeader *multipart.FileHeader, taskId int64) (*model.FileResponseDto, *model.ErrorResponse) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.uploadTaskImage.start")

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

	taskImageDto, errBuildTaskImageDto := mapper.BuildTaskImageDto(multipartFileHeader, config.Props.MinioBucket)
	if errBuildTaskImageDto != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-build-file", model.Exception),
			Message: errBuildTaskImageDto.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	taskImage := fs.FileRepo.SaveTaskImage(taskImageDto.TaskImage)

	taskTaskImage := &model.TaskTaskImage{
		TaskID:      taskId,
		TaskImageID: taskImage.Id,
	}

	errSave := fs.FileRepo.SaveTaskTaskImage(taskTaskImage)
	if errSave != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-save-task-image", model.Exception),
			Message: errSave.Error(),
			Code:    http.StatusForbidden,
		}
	}

	// Initialize Minio client
	minioClient, err := config.NewMinioClient()
	if err != nil {
		log.Fatalf("Failed to initialize Minio client: %v", err)
	}

	errUploadFileToMinio := util.UploadFileToMinio(ctx, taskImageDto.UniqueName, *multipartFile, multipartFileHeader.Size, minioClient)
	if errUploadFileToMinio != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-uplaod-to-minio", model.Exception),
			Message: errUploadFileToMinio.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	logger.Info("ActionLog.uploadTaskImage.end")
	return &model.FileResponseDto{taskImage.Id}, nil
}

func (fs *FileService) UploadTaskVideo(ctx context.Context, multipartFile *multipart.File, multipartFileHeader *multipart.FileHeader, taskId int64) (*model.FileResponseDto, *model.ErrorResponse) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.UploadTaskVideo.start")

	// Check if competitionId or multipartFile is nil
	if taskId <= 0 || multipartFileHeader == nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-continue", model.Exception),
			Message: "Task ID and file must not be null or empty",
			Code:    http.StatusBadRequest,
		}
	}

	// Check if file size exceeds limit 200MB
	maxSize, err := strconv.ParseInt(config.Props.TaskVideoMaxSize, 10, 64)
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

	taskVideoDto, errBuildTaskImageDto := mapper.BuildTaskVideoDto(multipartFileHeader, config.Props.MinioBucket)
	if errBuildTaskImageDto != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-build-file", model.Exception),
			Message: errBuildTaskImageDto.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	taskVideo := fs.FileRepo.SaveTaskVideo(taskVideoDto.TaskVideo)

	taskTaskVideo := &model.TaskTaskVideo{
		TaskID:      taskId,
		TaskVideoID: taskVideo.Id,
	}

	errSave := fs.FileRepo.SaveTaskTaskVideo(taskTaskVideo)
	if errSave != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-save-task-video", model.Exception),
			Message: errSave.Error(),
			Code:    http.StatusForbidden,
		}
	}

	// Initialize Minio client
	minioClient, err := config.NewMinioClient()
	if err != nil {
		log.Fatalf("Failed to initialize Minio client: %v", err)
	}

	errUploadFileToMinio := util.UploadFileToMinio(ctx, taskVideoDto.UniqueName, *multipartFile, multipartFileHeader.Size, minioClient)
	if errUploadFileToMinio != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-uplaod-to-minio", model.Exception),
			Message: errUploadFileToMinio.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	logger.Info("ActionLog.UploadTaskVideo.end")
	return &model.FileResponseDto{taskVideo.Id}, nil
}

func (fs *FileService) GetTaskImage(ctx context.Context, taskId int64, w http.ResponseWriter) *model.ErrorResponse {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GetTaskImage.start")

	taskTaskImage, errFindTaskTaskImage := fs.FileRepo.FindTaskTaskImageByTaskId(taskId)
	if errFindTaskTaskImage != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-find-task-task-image--file", model.Exception),
			Message: errFindTaskTaskImage.Error(),
			Code:    http.StatusForbidden,
		}
	}

	filePath := taskTaskImage.TaskImage.FilePath
	objectName := filePath[strings.LastIndex(filePath, "/")+1:]

	minioClient, errMinioClient := config.NewMinioClient()
	if errMinioClient != nil {
		log.Fatalf("Failed to initialize Minio client: %v", errMinioClient)
	}

	// Get object from MinIO
	object, err := minioClient.GetObject(context.Background(), config.Props.MinioBucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-get-file-from-minio", model.Exception),
			Message: err.Error(),
			Code:    http.StatusForbidden,
		}
	}
	defer object.Close()

	// Read a few bytes to detect content type
	buffer := make([]byte, 512) // Read first 512 bytes for MIME type detection
	_, err = object.Read(buffer)
	if err != nil && err != io.EOF {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-read-file-header", model.Exception),
			Message: err.Error(),
			Code:    http.StatusForbidden,
		}
	}

	// Reset object reader since we read some bytes
	object.Seek(0, io.SeekStart)

	// Detect MIME type
	contentType := http.DetectContentType(buffer)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "inline") // Remove filename to prevent download

	// Stream file to response
	_, err = io.Copy(w, object)
	if err != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-stream-file", model.Exception),
			Message: err.Error(),
			Code:    http.StatusForbidden,
		}
	}

	logger.Info("ActionLog.GetTaskImage.end")
	return nil
}

func (fs *FileService) StreamTaskVideo(r *http.Request, taskId int64, w http.ResponseWriter) *model.ErrorResponse {
	// Initialize Minio client
	ctx := r.Context()
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.StreamTaskVideo.start")

	taskTaskVideo, errFindTaskTaskVideo := fs.FileRepo.FindTaskTaskVideoByTaskId(taskId)
	if errFindTaskTaskVideo != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-find-task-task-video-file", model.Exception),
			Message: errFindTaskTaskVideo.Error(),
			Code:    http.StatusForbidden,
		}
	}

	filePath := taskTaskVideo.TaskVideo.FilePath
	objectName := filePath[strings.LastIndex(filePath, "/")+1:]

	minioClient, err := config.NewMinioClient()
	if err != nil {
		log.Fatalf("Failed to initialize Minio client: %v", err)
	}

	// Determine MIME type
	mimeType := "application/octet-stream"
	ext := strings.ToLower(filepath.Ext(objectName))
	if ext == ".mp4" {
		mimeType = "video/mp4"
	} else if ext == ".avi" {
		mimeType = "video/x-msvideo"
	} else if ext == ".mkv" {
		mimeType = "video/x-matroska"
	}

	// Get file size from MinIO
	objInfo, err := minioClient.StatObject(ctx, config.Props.MinioBucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.file-not-found", model.Exception),
			Message: err.Error(),
			Code:    http.StatusNotFound,
		}
	}
	fileSize := objInfo.Size

	// Set default range
	var start int64 = 0
	end := fileSize - 1

	// Parse Range Header
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" && strings.HasPrefix(rangeHeader, "bytes=") {
		ranges := strings.Split(rangeHeader[6:], "-")
		start, _ = strconv.ParseInt(ranges[0], 10, 64)
		if len(ranges) > 1 && ranges[1] != "" {
			end, _ = strconv.ParseInt(ranges[1], 10, 64)
		}
	}

	// Ensure range validity
	if start > end || start >= fileSize {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.invalid-range", model.Exception),
			Message: err.Error(),
			Code:    http.StatusRequestedRangeNotSatisfiable,
		}
	}

	contentLength := end - start + 1

	// Set response headers
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.WriteHeader(http.StatusPartialContent)

	// Stream file from MinIO
	opts := minio.GetObjectOptions{}
	opts.SetRange(start, end)

	obj, err := minioClient.GetObject(ctx, config.Props.MinioBucket, objectName, opts)
	if err != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-retrieve-file", model.Exception),
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	defer obj.Close()

	// Copy data to response
	buffer := make([]byte, 8192)
	for {
		n, err := obj.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("Error reading file:", err)
			return &model.ErrorResponse{
				Error:   fmt.Sprintf("%s.error-streaming-file", model.Exception),
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			}
		}
		_, writeErr := w.Write(buffer[:n])
		if writeErr != nil {
			log.Println("Client aborted connection:", writeErr)
			return &model.ErrorResponse{
				Error:   fmt.Sprintf("%s.client-aborted-connection", model.Exception),
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			}
		}
	}

	logger.Info("ActionLog.StreamTaskVideo.end")
	return nil
}
