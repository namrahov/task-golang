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

	router.HandleFunc(config.RootPath+"/users/login", h.authenticate).Methods("POST")
	router.HandleFunc(config.RootPath+"/users/register", h.register).Methods("POST")

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
		util.ErrorRespondWriterJSON(w, errRegister)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
