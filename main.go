package main

import (
	"compost/handlers"
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func main() {

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(pprof.New())

	// Initializing Data
	if err := handlers.InitDatasources(); err != nil {
		log.Println(err)
	}

	v1 := app.Group("/api")

	v1.Get("/account", handlers.Account)
	v1.Post("/login", handlers.Login)
	// v1.Get("/logout", handlers.logout)
	v1.Get("/note", handlers.GetTasks)
	v1.Post("/note", handlers.SaveOrEditTask)
	v1.Get("/note/:id", handlers.GetTask)
	v1.Post("/note/:id", handlers.DeleteTask)

	app.Static("/", "./static/public")

	app.Listen(":8080")
}
