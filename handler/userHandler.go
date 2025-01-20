package handler

import (
	"fmt"
	mid "github.com/go-chi/chi/middleware"
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

	router.Use(mid.Recoverer)
	//router.Use(middleware.RequestParamsMiddleware)

	h := &userHandler{
		UserService: &service.UserService{
			UserRepo: &repo.UserRepo{},
		},
	}

	router.HandleFunc("/users/register", h.register).Methods("POST")

	return router
}

func (h *userHandler) register(w http.ResponseWriter, r *http.Request) {
	var dto *model.UserRegistrationDto
	err := util.DecodeBody(w, r, &dto)
	if err != nil {
		return
	}

	fmt.Println(dto)
	errRegister := h.UserService.Register(r.Context(), dto)
	if errRegister != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}
