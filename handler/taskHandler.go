package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"task-golang/config"
	"task-golang/model"
	"task-golang/repo"
	"task-golang/service"
	"task-golang/util"
)

type taskHandler struct {
	TaskService service.ITaskService
}

func TaskHandler(router *mux.Router) *mux.Router {

	h := &taskHandler{
		TaskService: &service.TaskService{
			TaskRepo:  &repo.TaskRepo{},
			BoardRepo: &repo.BoardRepo{},
			UserUtil:  &util.UserUtil{},
		},
	}

	router.HandleFunc(config.RootPath+"/tasks/{boardId}", h.createTask).Methods("POST")

	return router
}

func (h *taskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	boardIdStr := mux.Vars(r)["id"]
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
