package main

import (
	"backend/configs"
	"fmt"

	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()
	db := configs.ConnectDB()

	jwt := middlewares.NewAuthMiddleware(configs.JWTSecret)
	app.Server.UseMiddleware(jwt)
	routes.UserRoute(app, jwt)
	routes.TodoRoutes(app, jwt)
	routes.MiscRoutes(app)
	
	app.Start()
}
