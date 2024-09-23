package server

import (
	Application "ZED-Magdy/Delivery-go/Application/Handlers"
	"ZED-Magdy/Delivery-go/Middlewares"
)

func (app *FiberServer) SetupRoutes() {
	userHandler := Application.NewUserHandler(app.DbService, app.ValidatorService)
	initInfoHandler := Application.NewInitInfoHandler(app.DbService)

	//Auth routes
	app.Post("/api/register", userHandler.CreateUser)
	app.Post("/api/login", userHandler.Login)

	//Authenticated routes
	authenticatedRoutes := app.Group("/api", Middlewares.AuthMiddleware)
	authenticatedRoutes.Get("/user", userHandler.GetCurrentUser)
	authenticatedRoutes.Get("/init-info", initInfoHandler.GetInitialInformation)
}
