package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	mails_service "odyssey/m/v2/internal/mail"
	"odyssey/m/v2/internal/token"
	"odyssey/m/v2/models/mail"
	"odyssey/m/v2/repositories"
	"odyssey/m/v2/usecases/mails"
	"strings"

	"gorm.io/gorm"
)

type MailsController struct {
	usecase      *mails.MailUseCase
	tokenService *token.TokenService
}

func NewMailController(db *gorm.DB, tokenService *token.TokenService, mailService *mails_service.MailService) *MailsController {
	mailRepo, err := repositories.NewMailRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	userRepo, err := repositories.NewUserRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	mailUseCase, err := mails.NewMailCase(mailRepo, userRepo, mailService)
	if err != nil {
		log.Fatal(err)
	}

	return &MailsController{usecase: mailUseCase, tokenService: tokenService}
}

func (c *MailsController) Create(w http.ResponseWriter, r *http.Request) {
	var mail mail.Mail
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&mail)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	jwtToken := strings.Split(r.Header["Authorization"][0], " ")[1]

	userId, err := c.tokenService.GetUserInfo(jwtToken)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	createdMail, err := c.usecase.Create(&mail, userId)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	jsonResp, _ := json.Marshal(createdMail)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
	w.WriteHeader(http.StatusOK)
}

func (c *MailsController) Get(w http.ResponseWriter, r *http.Request) {
	jwtToken := strings.Split(r.Header["Authorization"][0], " ")[1]

	userId, err := c.tokenService.GetUserInfo(jwtToken)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	result, err := c.usecase.Get(userId)
	if err != nil {
		c.errorResponse(w, "Bad request "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	jsonResp, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
	w.WriteHeader(http.StatusOK)
}

func (c *MailsController) errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
