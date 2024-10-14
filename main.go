package main

import (
	"mongoos/configs"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	configs.ConnectDB()
	app.Listen("127.0.0.1:6455")
}
