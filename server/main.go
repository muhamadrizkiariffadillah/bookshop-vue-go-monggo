package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/muhamadrizkiariffadillah/bookshop-vue-go-monggo/config"
)

func main() {
	config.LoadEnvVariable()

	port := config.GetEnvProperties("port")

	app := fiber.New()

	app.Use(recover.New())

	v1 := app.Group("/v1")
	v1.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(fiber.Map{"code": 200, "message": "ping-pong"})
	})

	err := app.Listen(":" + port)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
