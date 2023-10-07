package api

import (
	"github.com/gofiber/fiber/v2"

	"app/api/auth"
	"app/api/routes/userRoutes"
)

func testHandler(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func SetRoutes(app *fiber.App) {
	api := app.Group("/")
	userRoutes.SetUserRoutes(api)

	api.Get("/", auth.ValidateJWT, testHandler)
}

func SetupAPI(app *fiber.App) {
	SetRoutes(app)
}
