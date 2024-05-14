package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
    "github.com/rayhanmp/url-shortener-backend/db"
    "github.com/rayhanmp/url-shortener-backend/handlers"
	"github.com/redis/go-redis/v9"
	"context"
	"os"
)

func main() {
	ctx := context.Background()
    db.InitDB()
    redisURL := os.Getenv("REDIS_URL")
    opt, _ := redis.ParseURL(redisURL)
	
    rdb := redis.NewClient(opt)
    if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
    }

	log.Printf("Redis client in MAIN: %+v", rdb)

	app := fiber.New()
	app.Post("/shorten", func(c *fiber.Ctx) error {
		return handlers.ShortenURL(c, rdb)
	})

	app.Get("/:shortCode", func(c *fiber.Ctx) error {
		return handlers.RedirectURL(c, rdb)
	})

	log.Fatal(app.Listen(":8080"))
}