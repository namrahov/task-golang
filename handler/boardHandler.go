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

type boardHandler struct {
	BoardService service.IBoardService
}

func BoardHandler(router *mux.Router) *mux.Router {

	h := &boardHandler{
		BoardService: &service.BoardService{
			BoardRepo: &repo.BoardRepo{},
			UserRepo:  &repo.UserRepo{},
			UserUtil: &util.UserUtil{
				UserRepo: &repo.UserRepo{},
			},
		},
	}

	router.HandleFunc(config.RootPath+"/boards", h.create).Methods("POST")

	return router
}

// @Summary Create a new board
// @Description Creates a new board based on the provided data
// @Tags Boards
// @Accept  json
// @Produce  json
// @Param BoardRequestDto body model.BoardRequestDto true "Board Request Data"
// @Success 200 {string} string "Board created successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/boards [post]
// @Security     BearerAuth
func (h *boardHandler) create(w http.ResponseWriter, r *http.Request) {
	var dto *model.BoardRequestDto
	err := util.DecodeBody(w, r, &dto)
	if err != nil {
		return
	}

	errCreate := h.BoardService.CreateBoard(r.Context(), dto)
	if errCreate != nil {
		util.ErrorRespondWriterJSON(w, errCreate)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
