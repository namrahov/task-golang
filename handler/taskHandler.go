package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"task-golang/config"
	"task-golang/model"
	"task-golang/service"
	"task-golang/util"
)

type taskHandler struct {
	TaskService service.ITaskService
}

func TaskHandler(router *mux.Router, taskService *service.TaskService) *mux.Router {
	h := &taskHandler{
		taskService,
	}

	router.HandleFunc(config.RootPath+"/tasks/{boardId}", h.createTask).Methods("POST")

	return router
}

// @Summary Create a task
// @Description Create a new task for a specific board
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param boardId path int true "Board ID"
// @Param task body model.TaskRequestDto true "Task details"
// @Success 201 {string} string "Task created successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/tasks/{boardId} [post]
// @Security BearerAuth
func (h *taskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	boardIdStr := mux.Vars(r)["boardId"]
	boardId, err := strconv.ParseInt(boardIdStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dto *model.TaskRequestDto
	errGetDto := util.DecodeBody(w, r, &dto)
	if errGetDto != nil {
		return
	}

	errCreate := h.TaskService.CreateTask(r.Context(), dto, boardId)
	if errCreate != nil {
		util.ErrorRespondWriterJSON(w, errCreate)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
