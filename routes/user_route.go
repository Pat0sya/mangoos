package routes

import (
	"mongoos/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	// Root route to handle base URL
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the User API!")
	})

	// User-related routes
	app.Post("/user", controllers.CreateUser)
	app.Get("/user/:userId", controllers.GetAUser)
	app.Put("/user/:userId", controllers.EditAUser)
	app.Delete("/user/:userId", controllers.DeleteAUser)
	app.Get("/users", controllers.GetAllUsers)
}
