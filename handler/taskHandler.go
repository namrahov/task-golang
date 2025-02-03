package handler

import (
	"encoding/json"
	"fmt"
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
	router.HandleFunc(config.RootPath+"/tasks/{id:[0-9]+}", h.getTask).Methods("GET")
	router.HandleFunc(config.RootPath+"/tasks/page", h.getTasks).Methods("GET")

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

// @Summary Get task details
// @Description Retrieves the details of a specific task by ID
// @Tags Tasks
// @Accept json
// @Produce application/json
// @Param id path int true "Task ID"
// @Success 200 {object} model.TaskResponseDto "Task details retrieved successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request or task ID"
// @Failure 404 {object} model.ErrorResponse "Task not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/tasks/{id} [get]
// @Security BearerAuth
func (h *taskHandler) getTask(w http.ResponseWriter, r *http.Request) {
	// Parse the competition ID from the URL
	vars := mux.Vars(r)
	taskIDStr := vars["id"]
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, errGetTask := h.TaskService.GetTask(r.Context(), taskID)
	if errGetTask != nil {
		util.ErrorRespondWriterJSON(w, errGetTask)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// @Summary Get list of tasks
// @Description Retrieves a list of tasks with optional filtering by name, priority, and board ID, with pagination support
// @Tags Tasks
// @Accept json
// @Produce application/json
// @Param name query string false "Task name filter"
// @Param priority query string false "Task priority filter"
// @Param board_id query int false "Board ID to filter tasks"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param count query int false "Number of tasks per page (default: 10)"
// @Success 200 {object} model.TaskPageResponseDto "Tasks retrieved successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/tasks/page [get]
// @Security BearerAuth
func (h *taskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters
	name := r.URL.Query().Get("name")
	priority := r.URL.Query().Get("priority")
	boardIDStr := r.URL.Query().Get("board_id")
	pageStr := r.URL.Query().Get("page")
	countStr := r.URL.Query().Get("count")

	// Default pagination values
	page := 1
	count := 10
	var boardID int64

	// Convert board_id to int64 only if it's provided
	if boardIDStr != "" {
		var err error
		boardID, err = strconv.ParseInt(boardIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid board ID", http.StatusBadRequest)
			fmt.Println("Error parsing board_id:", err) // Log error
			return
		}
	}

	// Convert page and count to integers
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(countStr); err == nil && l > 0 {
		count = l
	}

	// Fetch tasks
	response, errGetTasks := h.TaskService.GetTasks(r.Context(), name, priority, boardID, page, count)
	if errGetTasks != nil {
		util.ErrorRespondWriterJSON(w, errGetTasks)
		return
	}

	// Encode and return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
