package app

import (
	"crypto/rsa"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"odyssey/m/v2/internal/mail"
	"odyssey/m/v2/internal/token"
)

type Application struct {
	db           *gorm.DB
	router       *http.Handler
	tokenService *token.TokenService
	mailService  *mail.MailService
}

func NewApplication() *Application {

	//dsn := "host=app-network.postgres user=postgres password=234492 dbname=currencies port=5432 sslmode=disable"
	dsn := "postgres://postgres:234492@postgres:5432/currencies?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	db.Logger.LogMode(4)

	mailService := mail.NewMailService(&mail.Config{Addr: "mail:2525"})

	return &Application{db: db, mailService: mailService}
}

func (app *Application) GetDb() *gorm.DB {
	return app.db
}

func (app *Application) GetRouter() *http.Handler {
	return app.router
}

func (app *Application) GetTokenService() *token.TokenService {
	return app.tokenService
}

func (app *Application) GetMailService() *mail.MailService {
	return app.mailService
}

func (app *Application) GetPublicKey() *rsa.PublicKey {
	return app.tokenService.PublicKey
}

func (app *Application) SetRouter(router *http.Handler) {
	app.router = router
}

func (app *Application) SetTokenService(service *token.TokenService) {
	app.tokenService = service
}
