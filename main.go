package main

import (
	"mongoos/configs"
	"mongoos/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes.UserRoute(app)
	configs.ConnectDB()
	app.Listen("127.0.0.1:6455")

}
