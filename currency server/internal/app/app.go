package app

import (
	"log"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"odyssey/m/v2/internal/mail"
	"odyssey/m/v2/internal/token"
)

type Application struct {
	db           *gorm.DB
	router       *mux.Router
	tokenService *token.TokenService
	mailService  *mail.MailService
}

func NewApplication() *Application {

	dsn := "host=localhost user=postgres password=234492 dbname=currencies port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	mailService := mail.NewMailService(&mail.Config{Addr: "localhost:2626"})

	return &Application{db: db, mailService: mailService}
}

func (app *Application) GetDb() *gorm.DB {
	return app.db
}

func (app *Application) GetRouter() *mux.Router {
	return app.router
}

func (app *Application) GetTokenService() *token.TokenService {
	return app.tokenService
}

func (app *Application) GetMailService() *mail.MailService {
	return app.mailService
}

func (app *Application) SetRouter(router *mux.Router) {
	app.router = router
}

func (app *Application) SetTokenService(service *token.TokenService) {
	app.tokenService = service
}
