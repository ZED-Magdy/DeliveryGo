package main

import (
	"ZED-Magdy/Delivery-go/Handlers"
	"ZED-Magdy/Delivery-go/Middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	db := initDb()
	trans, validate := initValidator()

	//Handlers
	userHandler := Handlers.NewUserHandler(*db, validate, trans)
	initInfoHandler := Handlers.NewInitInfoHandler(*db)

	//Auth routes
	app.Post("/api/register", userHandler.CreateUser)
	app.Post("/api/login", userHandler.Login)

	//Authenticated routes
	authenticatedRoutes := app.Group("/api", Middlewares.AuthMiddleware)
	authenticatedRoutes.Get("/user", userHandler.GetCurrentUser)
	authenticatedRoutes.Get("/init-info", initInfoHandler.GetInitialInformation)
}
