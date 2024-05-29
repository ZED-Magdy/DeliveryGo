package services

import (
	"ZED-Magdy/Delivery-go/Models"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)


type AuthService struct {
	request *http.Request
	db	  gorm.DB
}

func NewAuthService(request *http.Request, db gorm.DB) *AuthService {
	return &AuthService{
		request: request,
		db: db,
	}
}

func (a *AuthService) User() (*Models.User, error) {
	authHeader := strings.Split(a.request.Header.Get("Authorization"), "Bearer ")
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