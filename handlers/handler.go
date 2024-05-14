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


	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	shortCode := utils.StringWithCharset(shortCodeLength, charset)
	url := models.URL{
		OriginalURL: request.URL,
		ShortCode:   shortCode,
	}

	result := db.DB.Create(&url)
		if result.Error != nil {
			log.Printf("Error creating URL record: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": result.Error.Error(),
			})
		}

		err := rdb.Set(ctx, shortCode, request.URL, 0).Err()
		if err != nil {
			panic(err)
		}
		log.Printf("%v", err)
		log.Printf("Redis client in SHORTENURL: %+v", rdb)
		log.Printf("Caching success")
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

	result := db.DB.First(&url, "short_code = ?", shortCode)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found",
		})
	}

	return c.Redirect(url.OriginalURL, fiber.StatusFound)
}
