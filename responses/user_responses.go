package responses

import "github.com/gofiber/fiber"

type UserResonce struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}
