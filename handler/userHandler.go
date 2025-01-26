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
	router.HandleFunc(config.RootPath+"/users/demo", h.demo).Methods("POST")

	return router
}

func (h *userHandler) demo(w http.ResponseWriter, r *http.Request) {
	//idStr := mux.Vars(r)["id"]
	//id, err := strconv.ParseInt(idStr, 10, 64)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	fmt.Println("Demo isledi senik")
}

// authenticate handles user authentication and generates a JWT token.
//
// @Summary      Authenticate user
// @Description  Authenticates a user by validating their credentials and returns a JWT token upon success.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        authRequestDto body model.AuthRequestDto true "authRequestDto"
// @Success      200 {object} model.JwtToken "JWT Token response"
// @Failure      400 {object} model.ErrorResponse "Bad Request"
// @Failure      401 {object} model.ErrorResponse "Unauthorized"
// @Failure      500 {object} model.ErrorResponse "Internal Server Error"
// @Router       /v1/users/login [post]
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

// register handles user registration.
// @Summary Register a new user
// @Description Registers a new user with the provided information.
// @Tags User
// @Accept json
// @Produce json
// @Param userRegistrationDto body model.UserRegistrationDto true "userRegistrationDto"
// @Success 201 {string} string "User successfully registered"
// @Failure 400 {object} model.ErrorResponse "Invalid request data"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/users/register [post]
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

// active activates a user account.
// @Summary Activate user account
// @Description Activates a user account using the provided activation token.
// @Tags User
// @Accept json
// @Produce json
// @Param token query string true "Activation token"
// @Success 204 {string} string "User successfully activated"
// @Failure 400 {string} string "Token is required"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /v1/users/activate [get]
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
