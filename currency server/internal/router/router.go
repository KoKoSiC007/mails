package router

import (
	"net/http"
	"odyssey/m/v2/controllers"
	"odyssey/m/v2/internal/app"
	"odyssey/m/v2/internal/logger"
	"odyssey/m/v2/usecases/users"
	"strings"

	"github.com/gorilla/mux"
)

type UserController struct {
	usecase *users.UserUseCase
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(app *app.Application) {
	router := mux.NewRouter()

	userController := controllers.NewUserController(app.GetDb(), app.GetTokenService())
	currencyController := controllers.NewCurrencyController(app.GetDb())
	mailsController := controllers.NewMailController(app.GetDb(), app.GetTokenService(), app.GetMailService())

	authorization := app.GetTokenService().ValidateToken

	var routes = Routes{
		Route{
			"GetReport",
			strings.ToUpper("Get"),
			"/api/v1/currencies",
			currencyController.Get,
		},

		Route{
			"GetMails",
			strings.ToUpper("Get"),
			"/api/v1/mail",
			mailsController.Get,
		},

		Route{
			"SendMail",
			strings.ToUpper("Post"),
			"/api/v1/mail",
			mailsController.Create,
		},

		Route{
			"CreateUser",
			strings.ToUpper("Post"),
			"/api/v1/user",
			userController.Create,
		},

		Route{
			"LoginUser",
			strings.ToUpper("Get"),
			"/api/v1/user/login",
			userController.Auth,
		},
	}

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)
		if route.Name == "GetReport" || route.Name == "SendMail" || route.Name == "GetMails" {
			handler = authorization(handler)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	app.SetRouter(router)
}
