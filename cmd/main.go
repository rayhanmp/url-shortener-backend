package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rayhanmp/url-shortener-backend/db"
)

func main() {
	app := fiber.New()

	db.CreateClient()

	setupRoutes(app)

	app.Listen(":3000")
}