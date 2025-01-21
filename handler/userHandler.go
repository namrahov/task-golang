package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"task-golang/model"
	"task-golang/repo"
	"task-golang/service"
	"task-golang/util"
)

type userHandler struct {
	UserService service.IUserService
}

func UserHandler(router *mux.Router) *mux.Router {

	h := &userHandler{
		UserService: &service.UserService{
			UserRepo:        &repo.UserRepo{},
			PasswordChecker: &util.PasswordChecker{},
		},
	}

	router.HandleFunc("/users/login", h.authenticate).Methods("POST")
	router.HandleFunc("/users/register", h.register).Methods("POST")

	return router
}

func (h *userHandler) authenticate(w http.ResponseWriter, r *http.Request) {

}

func (h *userHandler) register(w http.ResponseWriter, r *http.Request) {
	var dto *model.UserRegistrationDto
	err := util.DecodeBody(w, r, &dto)
	if err != nil {
		return
	}

	errRegister := h.UserService.Register(r.Context(), dto)
	if errRegister != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(errRegister.Code)
		json.NewEncoder(w).Encode(errRegister)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
