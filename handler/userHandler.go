package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"task-golang/config"
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
			TokenRepo:       &repo.TokenRepo{},
			PasswordChecker: &util.PasswordChecker{},
			TokenUtil: &util.TokenUtil{
				TokenRepo: &repo.TokenRepo{},
			},
		},
	}

	router.HandleFunc(config.RootPath+"/users/login", h.authenticate).Methods("POST")
	router.HandleFunc(config.RootPath+"/users/register", h.register).Methods("POST")
	router.HandleFunc(config.RootPath+"/users/active", h.active).Methods("GET")

	return router
}

func (h *userHandler) authenticate(w http.ResponseWriter, r *http.Request) {
	var dto *model.AuthRequestDto
	err := util.DecodeBody(w, r, &dto)
	if err != nil {
		return
	}

	jwtToken, errLogin := h.UserService.Authenticate(r.Context(), dto)
	if errLogin != nil {
		util.ErrorRespondWriterJSON(w, errLogin)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jwtToken)
}

func (h *userHandler) register(w http.ResponseWriter, r *http.Request) {
	var dto *model.UserRegistrationDto
	err := util.DecodeBody(w, r, &dto)
	if err != nil {
		return
	}

	errRegister := h.UserService.Register(r.Context(), dto)
	if errRegister != nil {
		util.ErrorRespondWriterJSON(w, errRegister)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *userHandler) active(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}
	fmt.Println("token=", token)

	h.UserService.Active(r.Context(), token)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
