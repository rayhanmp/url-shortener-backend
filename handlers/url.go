package handlers

import (	
	"github.com/gofiber/fiber/v2"
	"github.com/rayhanmp/url-shortener-backend/models"
	"github.com/rayhanmp/url-shortener-backend/db"
)


func Home (c *fiber.Ctx) error {
	return c.SendString("Hello, Captain!")
}

func ShortenURL (c *fiber.Ctx) error {
	url := &models.URL{}

	if err := c.BodyParser(url); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	err := db.Supabase.DB.From("urls").Insert(url).Execute(&url)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Cannot insert URL",
		})
	}

	return c.JSON(url)
}