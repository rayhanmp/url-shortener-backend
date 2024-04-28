package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rayhanmp/url-shortener-backend/handlers"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Post("/shorten", handlers.ShortenURL)
}