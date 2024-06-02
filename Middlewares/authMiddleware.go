package Middlewares

import (
	services "ZED-Magdy/Delivery-go/Services"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	if c.Get("Authorization") == "" {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"message": "Unauthorized"})
	}

	authHeader := strings.Split(c.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"message": "Unauthorized"})
	}

	_, err := services.VerifyJwtToken(authHeader[1])
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"message": "Unauthorized"})
	}

	return c.Next()

}
