package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"task-golang/config"
	"task-golang/model"
	"task-golang/service"
	"task-golang/util"
)

type boardHandler struct {
	BoardService service.IBoardService
}

func BoardHandler(router *mux.Router, boardService *service.BoardService) *mux.Router {
	h := &boardHandler{
		boardService,
	}

	router.HandleFunc(config.RootPath+"/boards", h.create).Methods("POST")
	router.HandleFunc(config.RootPath+"/boards/{id}/access", h.giveAccess).Methods("POST")
	router.HandleFunc(config.RootPath+"/boards/{userId}", h.getUserBoards).Methods("GET")

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

// @Summary Give access to a board
// @Description Assign access to a specific board for a user
// @Tags Boards
// @Accept  json
// @Produce  json
// @Param id path int true "Board ID"
// @Param userId query int true "User ID"
// @Success 200 {string} string "Access granted successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/boards/{id}/access [post]
// @Security     BearerAuth
func (h *boardHandler) giveAccess(w http.ResponseWriter, r *http.Request) {
	boardIdStr := mux.Vars(r)["id"]
	boardId, err := strconv.ParseInt(boardIdStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userIdStr := r.URL.Query().Get("userId")
	// Convert to int64
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid boardId", http.StatusBadRequest)
		return
	}

	errCreate := h.BoardService.GiveAccessToBoard(r.Context(), userId, boardId)
	if errCreate != nil {
		util.ErrorRespondWriterJSON(w, errCreate)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// @Summary Get user boards
// @Description Retrieve all boards associated with a specific user
// @Tags Boards
// @Accept  json
// @Produce  json
// @Param userId path int true "User ID"
// @Success 200 {array} model.Board "List of user boards"
// @Failure 400 {object} model.ErrorResponse "Invalid request"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/boards/{userId} [get]
// @Security BearerAuth
func (h *boardHandler) getUserBoards(w http.ResponseWriter, r *http.Request) {
	userIdStr := mux.Vars(r)["userId"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersBoards, errCreate := h.BoardService.GetUserBoards(r.Context(), userId)
	if errCreate != nil {
		util.ErrorRespondWriterJSON(w, errCreate)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usersBoards)
}
