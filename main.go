package main

import (
	"ZED-Magdy/Delivery-go/infrastructure/server"
	"log"

	"github.com/joho/godotenv"
	// "github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	app := server.New()

	app.SetupRoutes()

	app.Listen(":8080")
}

// func initCache(ctx context.Context) *cache.Cache[string] {
// 	redisStore := redis_store.NewRedis(redis.NewClient(&redis.Options{
// 		Addr: "127.0.0.1:6379",
// 	}))

// 	cacheManager := cache.New[string](redisStore)

// 	return cacheManager
// }
