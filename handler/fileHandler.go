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
		http.Error(w, "Invalid competition ID", http.StatusBadRequest)
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
