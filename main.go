package main

import (
	"fiber-proj1/handlers"
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Initializing Data
	if err := handlers.InitDatasources(); err != nil {
		log.Println(err)
	}
	v1 := app.Group("/api")

	v1.Post("/login", handlers.Login)
	// v1.Get("/logout", handlers.logout)
	v1.Get("/task", handlers.GetTasks)
	v1.Post("/task", handlers.SaveOrEditTask)
	v1.Get("/task/:id", handlers.GetTask)
	v1.Post("/task/:id", handlers.DeleteTask)

	app.Listen(":8080")
}
