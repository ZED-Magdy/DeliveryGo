package services

import (
	"ZED-Magdy/Delivery-go/Models"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)


type AuthService struct {
	ctx *fiber.Ctx
	db	  gorm.DB
}

func NewAuthService(ctx *fiber.Ctx, db gorm.DB) *AuthService {
	return &AuthService{
		ctx: ctx,
		db: db,
	}
}

func (a *AuthService) User() (*Models.User, error) {
	authHeader := strings.Split(a.ctx.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return nil, errors.New("Unauthorized")
	}

	token := authHeader[1]
	tokenObj, err := VerifyJwtToken(token)

	if err != nil {
		return nil, err
	}

	user := Models.User{}
	err = a.db.First(&user, tokenObj.Claims.(jwt.MapClaims)["subject"].(float64)).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}