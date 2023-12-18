package routes

import (
	"backend/controllers"
	"fmt"

	"gofr.dev/pkg/gofr"
)

func UserRoute(app *gofr.App, authMiddleware gofr.Handler) {
	fmt.Print(app)
	app.Post("/user", controllers.CreateUser)
	app.Post("/login", controllers.LoginUser)
	app.Get("/users/:userId", controllers.GetUser)
	app.Get("/user", authMiddleware, controllers.GetProfile)
	app.Put("/user", authMiddleware, controllers.EditPassword)
}
