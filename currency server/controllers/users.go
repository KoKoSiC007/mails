package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"odyssey/m/v2/internal/token"
	"odyssey/m/v2/models/user"
	"odyssey/m/v2/repositories"
	"odyssey/m/v2/usecases/users"

	"gorm.io/gorm"
)

type UserController struct {
	usecase      *users.UserUseCase
	tokenService *token.TokenService
}

func NewUserController(db *gorm.DB, tokenService *token.TokenService) *UserController {
	userRepo, err := repositories.NewUserRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	userUseCase, err := users.NewUserCase(userRepo)
	if err != nil {
		log.Fatal(err)
	}

	return &UserController{usecase: userUseCase, tokenService: tokenService}
}

func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var user user.NewUser
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&user)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	createdUser, err := c.usecase.Create(&user)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	jsonResp, _ := json.Marshal(createdUser)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
	w.WriteHeader(http.StatusCreated)
}
func (c *UserController) Auth(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")
	pass := r.URL.Query().Get("password")
	if len(login) == 0 || len(pass) == 0 {
		c.errorResponse(w, "login or pass is missing", http.StatusUnprocessableEntity)
		return
	}

	user, err := c.usecase.Auth(login, pass)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if user == nil {
		c.errorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	token, err := token.GenerateToken(&user.Id, c.tokenService.PrivateKey)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusCreated)
}

func (c *UserController) errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
