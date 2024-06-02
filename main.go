package main

import (
	"ZED-Magdy/Delivery-go/Handlers"
	"ZED-Magdy/Delivery-go/Middlewares"
	"encoding/json"
	"strings"
	"ZED-Magdy/Delivery-go/Models"
	"log"
	"os"
	// "github.com/eko/gocache/lib/v4/cache"
	// redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	// "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint:  "/health",
		ReadinessEndpoint: "/ready",
	}))

	app.Use(idempotency.New())
	app.Use(recover.New())
	db := initDb()
	trans, validate := initValidator()
	userHandler := Handlers.NewUserHandler(*db, validate, trans)
	app.Post("/api/register", userHandler.CreateUser)
	app.Post("/api/login", userHandler.Login)
	initInfoHandler := Handlers.NewInitInfoHandler(*db)
	authenticatedRoutes := app.Group("/api", Middlewares.AuthMiddleware)
	authenticatedRoutes.Get("/user", userHandler.GetCurrentUser)
	authenticatedRoutes.Get("/init-info", initInfoHandler.GetInitialInformation)
	app.Listen(":8000")
}

func initValidator() (ut.Translator, *validator.Validate) {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
	return trans, validate
}

func initDb() *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Models.User{}, &Models.PriceList{}, &Models.Region{})
	seed(db)
	return db
}
func seed(db *gorm.DB) {
	user := &Models.User{}
	db.First(user)
	if user.ID == 0 {
		db.Create(&Models.PriceList{Name: "Standard", KmCost: 0.5, CancellationCost: 0.5})
		db.Exec("INSERT INTO regions (name, geofence, price_list_id) VALUES (?, ST_GeomFromText(?), ?)", "Cairo", "POLYGON((30.0444 31.2357, 30.0444 30.0444, 31.2357 30.0444, 31.2357 31.2357, 30.0444 31.2357))", 1)
	}

}

// func initCache(ctx context.Context) *cache.Cache[string] {
// 	redisStore := redis_store.NewRedis(redis.NewClient(&redis.Options{
// 		Addr: "127.0.0.1:6379",
// 	}))

// 	cacheManager := cache.New[string](redisStore)

// 	return cacheManager
// }
