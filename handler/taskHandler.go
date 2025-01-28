package handler

import (
	"github.com/gorilla/mux"
	"net/http"
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
			TaskRepo: &repo.TaskRepo{},
		},
	}

	router.HandleFunc(config.RootPath+"/tasks", h.createTask).Methods("POST")

	return router
}

func (h *taskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var dto *model.BoardRequestDto
	err := util.DecodeBody(w, r, &dto)
	if err != nil {
		return
	}

	//errCreate := h.BoardService.CreateBoard(r.Context(), dto)
	//if errCreate != nil {
	//	util.ErrorRespondWriterJSON(w, errCreate)
	//	return
	//}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
