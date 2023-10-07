package server

import (
	"app/api"
	"app/database"
	"fmt"
	"log"

	"app/config"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Start() {
	// Initalize Database Or Die
	database.InitDB()
	// Create New App With Faster JSON Encoder
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Set Routes & Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	api.SetupAPI(app)

	APP_PORT := ":" + fmt.Sprintf("%d", config.APP_PORT)
	// Start API
	fmt.Printf("\nStarting app at http://localhost%s\n", APP_PORT)
	log.Fatal(app.Listen(APP_PORT))
}
