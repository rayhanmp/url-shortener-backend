package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rayhanmp/url-shortener-backend/db"
	"github.com/rayhanmp/url-shortener-backend/models"
	"github.com/rayhanmp/url-shortener-backend/utils"
	"github.com/redis/go-redis/v9"
)

const (
	shortCodeLength = 6
	charset         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var ctx = context.Background()

func ShortenURL(c *fiber.Ctx, rdb *redis.Client) error {
	var request struct {
		URL string `json:"url"`
	}

	// Parse the request body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate a random short code
	shortCode := utils.StringWithCharset(shortCodeLength, charset)
	url := models.URL{
		OriginalURL: request.URL,
		ShortCode:   shortCode,
	}

	// Save the URL in the database
	result := db.DB.Create(&url)
	if result.Error != nil {
		log.Printf("Error creating URL record: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}

	// Cache the URL in Redis
	err := rdb.Set(ctx, shortCode, request.URL, 0).Err()
	if err != nil {
		panic(err)
	}

	response := map[string]string{"shortCode": shortCode}

	return c.JSON(response)
}

func RedirectURL(c *fiber.Ctx, rdb *redis.Client) error {
	shortCode := c.Params("shortCode")
	var url models.URL

	// Check if the URL is cached in Redis
    originalURL, err := rdb.Get(ctx, shortCode).Result()
    if err == nil {
		log.Printf("Cache hit: %v", originalURL)
        return c.Redirect(originalURL, fiber.StatusFound)
    }

	// If the URL is not cached in Redis, fetch it from the database
	result := db.DB.First(&url, "short_code = ?", shortCode)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found",
		})
	}

	// Cache the URL in Redis
	errRDB := rdb.Set(ctx, shortCode, url.OriginalURL, 0).Err()
	if errRDB != nil {
		panic(errRDB)
	}

	return c.Redirect(url.OriginalURL, fiber.StatusFound)
}
