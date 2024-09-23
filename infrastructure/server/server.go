package server

import (
	"ZED-Magdy/Delivery-go/Models"
	"ZED-Magdy/Delivery-go/infrastructure/database"
	"ZED-Magdy/Delivery-go/infrastructure/validator"
	"encoding/json"
	"strings"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

type FiberServer struct {
	*fiber.App
	DbService        *database.Service
	ValidatorService validator.Service
}

func New() *FiberServer {
	app := fiber.New(fiber.Config{
		ServerHeader: "ZadDelivery",
		AppName:      "ZadDelivery",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
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
		Level: compress.LevelBestSpeed,
	}))

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessEndpoint:  "/health",
		ReadinessEndpoint: "/ready",
	}))

	app.Use(idempotency.New())
	app.Use(recover.New())
	app.Use(swagger.New())
	dbService := database.New()
	seed(dbService.Db)

	return &FiberServer{
		App:              app,
		DbService:        dbService,
		ValidatorService: validator.New(),
	}
}
func seed(db *gorm.DB) {
	user := &Models.User{}
	db.First(user)
	if user.ID == 0 {
		db.Create(&Models.PriceList{Name: "Standard", KmCost: 0.5, CancellationCost: 0.5})
		db.Exec("INSERT INTO regions (name, geofence, price_list_id) VALUES (?, ST_GeomFromText(?), ?)", "Cairo", "POLYGON((30.0444 31.2357, 30.0444 30.0444, 31.2357 30.0444, 31.2357 31.2357, 30.0444 31.2357))", 1)
	}

}
