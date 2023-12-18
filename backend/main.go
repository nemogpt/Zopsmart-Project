package main

import (
	"backend/configs"
	"fmt"

	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()
	db := configs.ConnectDB()
	app.GET("/greet", func(ctx *gofr.Context) (interface{}, error) {
		fmt.Print(db)
		return "get api", nil
	})
	app.Post("/todo", func(ctx *gofr.Context) (interface{}, error) {
		fmt.Print(db)
		return "post api", nil
	})
	app.Get("/todo/:todoId", func(ctx *gofr.Context) (interface{}, error) {
		fmt.Print(db)
		return "todo with id", {todoId}
	})
	app.Put("/todo/:todoId",func(ctx *gofr.Context) (interface{}, error) {
		fmt.Print(db)
		return "todo with id", {todoId}
	})
	app.Delete("/todo/:todoId",func(ctx *gofr.Context) (interface{}, error) {
		fmt.Print(db)
		return "todo with id", {todoId}
	})
	app.Get("/todos", func(ctx *gofr.Context) (interface{}, error) {
		fmt.Print(db)
		return "Hello World!", nil
	})

	// jwt := middlewares.NewAuthMiddleware(configs.JWTSecret())

	routes.UserRoute(app, jwt)
	// routes.TodoRoutes(app, jwt)
	// routes.MiscRoutes(app)
	app.Start()
}
