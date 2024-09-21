package Handlers

import (
	"ZED-Magdy/Delivery-go/Dtos"
	"ZED-Magdy/Delivery-go/Models"
	services "ZED-Magdy/Delivery-go/Services"
	"net/http"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db        gorm.DB
	validator *validator.Validate
	trans     ut.Translator
}

func NewUserHandler(db gorm.DB, validator *validator.Validate, trans ut.Translator) *UserHandler {
	return &UserHandler{
		db:        db,
		validator: validator,
		trans:     trans,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	request := new(Dtos.CreateUserDto)
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "Invalid request payload"})
	}
	validationResult := services.Validate(h.validator, h.trans, request)
	if validationResult != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(validationResult)
	}
	hashedPassword, err := h.hashPassword(request.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"message": "Server Error"})
	}

	user := Models.User{
		Name:     request.Name,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: hashedPassword,
	}

	err = h.db.Create(&user).Error
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // MySQL code for duplicate entry
				return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "User already exists"})
			default:
				return c.Status(http.StatusInternalServerError).JSON(map[string]string{"message": "Error creating user"})

			}
		}
	}
	token, err := services.CreateJwtToken(user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"message": "Error creating a token"})
	}
	userDto := Dtos.UserDto{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Token: &token,
	}
	return c.Status(http.StatusCreated).JSON(userDto)

}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	request := new(Dtos.LoginDto)
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "Invalid request payload"})
	}
	validationResult := services.Validate(h.validator, h.trans, request)
	if validationResult != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(validationResult)
	}
	user := Models.User{}
	h.db.Where("phone = ?", request.Phone).First(&user)
	if user.ID == 0 {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"message": "Invalid credentials"})
	}

	if !h.checkPasswordHash(request.Password, user.Password) {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"message": "Invalid credentials"})
	}

	token, err := services.CreateJwtToken(user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"message": "Error creating a token"})
	}
	userDto := Dtos.UserDto{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Token: &token,
	}
	return c.Status(http.StatusOK).JSON(userDto)
}

func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	authService := services.NewAuthService(c, h.db)
	user, err := authService.User()
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"message": err.Error()})
	}
	token := strings.Split(c.Get("Authorization"), "Bearer ")[1]
	userDto := Dtos.UserDto{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Token: &token,
	}

	return c.Status(http.StatusOK).JSON(userDto)
}

func (*UserHandler) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return string(bytes), err
}

func (*UserHandler) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
