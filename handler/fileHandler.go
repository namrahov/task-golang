package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"task-golang/config"
	"task-golang/service"
	"task-golang/util"
)

type fileHandler struct {
	FileService service.IFileService
}

func FileHandler(router *mux.Router, fileService *service.FileService) *mux.Router {
	h := &fileHandler{
		fileService,
	}

	router.HandleFunc(config.RootPath+"/files/upload/attachment/{taskId}", h.uploadAttachmentFile).Methods("POST")
	router.HandleFunc(config.RootPath+"/files/delete/attachment/{attachmentFileId}", h.deleteAttachmentFile).Methods("DELETE")
	router.HandleFunc(config.RootPath+"/files/download/attachment/{attachmentFileId}", h.downloadAttachmentFile).Methods("GET")
	router.HandleFunc(config.RootPath+"/files/upload/task-image/{taskId}", h.uploadTaskImage).Methods("POST")
	router.HandleFunc(config.RootPath+"/files/get/task-image/{taskId}", h.getTaskImage).Methods("GET")
	router.HandleFunc(config.RootPath+"/files/stream/task-video/{taskVideoId}", h.streamTaskVideo).Methods("GET")
	router.HandleFunc(config.RootPath+"/files/upload/task-video/{taskId}", h.uploadTaskVideo).Methods("POST")

	return router
}

// @Summary Upload an attachment file
// @Description Uploads a file as an attachment for a specific task
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param taskId path int true "Task ID"
// @Param file formData file true "File to upload"
// @Success 201 {object} model.FileResponseDto "File uploaded successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request or file"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/files/upload/attachment/{taskId} [post]
// @Security BearerAuth
func (h *fileHandler) uploadAttachmentFile(w http.ResponseWriter, r *http.Request) {
	// Parse the competition ID from the URL
	vars := mux.Vars(r)
	taskIDStr := vars["taskId"]
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, multipartFileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Log the file name and size (optional)
	fmt.Printf("Uploaded File: %+v\n", multipartFileHeader.Filename)
	fmt.Printf("File Size: %+v\n", multipartFileHeader.Size)
	fmt.Printf("MIME Header: %+v\n", multipartFileHeader.Header)

	// Call the file service to handle the file upload
	response, errUpload := h.FileService.UploadAttachmentFile(r.Context(), &file, multipartFileHeader, taskID)
	if errUpload != nil {
		util.ErrorRespondWriterJSON(w, errUpload)
		return
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// @Summary Delete an attachment file
// @Description Deletes an attachment file associated with a specific task
// @Tags Files
// @Param attachmentFileId path int true "Attachment File ID"
// @Success 204 "File deleted successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request or file ID"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/files/delete/attachment/{attachmentFileId} [delete]
// @Security BearerAuth
func (h *fileHandler) deleteAttachmentFile(w http.ResponseWriter, r *http.Request) {
	// Parse the competition ID from the URL
	vars := mux.Vars(r)
	attachmentFileIdStr := vars["attachmentFileId"]
	attachmentFileId, err := strconv.ParseInt(attachmentFileIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid attachment file ID", http.StatusBadRequest)
		return
	}

	errDeleteAttachmentFile := h.FileService.DeleteAttachmentFile(r.Context(), attachmentFileId)
	if errDeleteAttachmentFile != nil {
		util.ErrorRespondWriterJSON(w, errDeleteAttachmentFile)
		return
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Download an attachment file
// @Description Downloads an attachment file associated with a specific task
// @Tags Files
// @Param attachmentFileId path int true "Attachment File ID"
// @Success 200 "File downloaded successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request or file ID"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/files/download/attachment/{attachmentFileId} [get]
// @Security BearerAuth
func (h *fileHandler) downloadAttachmentFile(w http.ResponseWriter, r *http.Request) {
	// Parse the competition ID from the URL
	vars := mux.Vars(r)
	attachmentFileIdStr := vars["attachmentFileId"]
	attachmentFileId, err := strconv.ParseInt(attachmentFileIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid attachment file ID", http.StatusBadRequest)
		return
	}

	errDeleteAttachmentFile := h.FileService.DownloadAttachmentFile(r.Context(), attachmentFileId, w)
	if errDeleteAttachmentFile != nil {
		util.ErrorRespondWriterJSON(w, errDeleteAttachmentFile)
		return
	}
}

// @Summary Upload an image for a task
// @Description Uploads an image file associated with a specific task
// @Tags Files
// @Accept multipart/form-data
// @Produce application/json
// @Param taskId path int true "Task ID"
// @Param file formData file true "Image file to upload"
// @Success 201 {object} model.FileResponseDto "File uploaded successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request or task ID"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/files/upload/task-image/{taskId} [post]
// @Security BearerAuth
func (h *fileHandler) uploadTaskImage(w http.ResponseWriter, r *http.Request) {
	// Parse the competition ID from the URL
	vars := mux.Vars(r)
	taskIDStr := vars["taskId"]
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, multipartFileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Log the file name and size (optional)
	fmt.Printf("Uploaded File: %+v\n", multipartFileHeader.Filename)
	fmt.Printf("File Size: %+v\n", multipartFileHeader.Size)
	fmt.Printf("MIME Header: %+v\n", multipartFileHeader.Header)

	// Call the file service to handle the file upload
	response, errUpload := h.FileService.UploadTaskImage(r.Context(), &file, multipartFileHeader, taskID)
	if errUpload != nil {
		util.ErrorRespondWriterJSON(w, errUpload)
		return
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// @Summary Get task image
// @Description Retrieves an image file associated with a specific task
// @Tags Files
// @Accept json
// @Produce image/png, image/jpeg, image/webp
// @Param taskId path int true "Task ID"
// @Success 200 {file} binary "Image file retrieved successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request or task ID"
// @Failure 404 {object} model.ErrorResponse "Image not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/files/get/task-image/{taskId} [get]
// @Security BearerAuth
func (h *fileHandler) getTaskImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIDStr := vars["taskId"]
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	errDeleteAttachmentFile := h.FileService.GetTaskImage(r.Context(), taskID, w)
	if errDeleteAttachmentFile != nil {
		util.ErrorRespondWriterJSON(w, errDeleteAttachmentFile)
		return
	}
}

// @Summary Stream task video
// @Description Streams a video file associated with a specific task
// @Tags Files
// @Accept json
// @Produce video/mp4, video/webm, video/ogg
// @Param taskVideoId path int true "Task Video ID"
// @Success 200 {file} binary "Video file streamed successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request or task ID"
// @Failure 404 {object} model.ErrorResponse "Video not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/files/stream/task-video/{taskVideoId} [get]
// @Security BearerAuth
func (h *fileHandler) streamTaskVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskVideoIdStr := vars["taskVideoId"]
	taskVideoId, err := strconv.ParseInt(taskVideoIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	errDeleteAttachmentFile := h.FileService.StreamTaskVideo(r, taskVideoId, w)
	if errDeleteAttachmentFile != nil {
		util.ErrorRespondWriterJSON(w, errDeleteAttachmentFile)
		return
	}
}

// @Summary Upload a video for a task
// @Description Uploads a video file associated with a specific task
// @Tags Files
// @Accept multipart/form-data
// @Produce application/json
// @Param taskId path int true "Task ID"
// @Param file formData file true "Video file to upload"
// @Success 201 {object} model.FileResponseDto "File uploaded successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request or task ID"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/files/upload/task-video/{taskId} [post]
// @Security BearerAuth
func (h *fileHandler) uploadTaskVideo(w http.ResponseWriter, r *http.Request) {
	// Parse the competition ID from the URL
	vars := mux.Vars(r)
	taskIDStr := vars["taskId"]
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, multipartFileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Log the file name and size (optional)
	fmt.Printf("Uploaded File: %+v\n", multipartFileHeader.Filename)
	fmt.Printf("File Size: %+v\n", multipartFileHeader.Size)
	fmt.Printf("MIME Header: %+v\n", multipartFileHeader.Header)

	// Call the file service to handle the file upload
	response, errUpload := h.FileService.UploadTaskVideo(r.Context(), &file, multipartFileHeader, taskID)
	if errUpload != nil {
		util.ErrorRespondWriterJSON(w, errUpload)
		return
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

//add these to getTask
//dovnload excel of tasks
