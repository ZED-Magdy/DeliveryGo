package Application

import (
	"ZED-Magdy/Delivery-go/Application/Dtos"
	"ZED-Magdy/Delivery-go/Models"
	services "ZED-Magdy/Delivery-go/Services"
	"ZED-Magdy/Delivery-go/infrastructure/database"
	"ZED-Magdy/Delivery-go/infrastructure/validator"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	dbService        *database.Service
	validatorService validator.Service
}

func NewUserHandler(dbService *database.Service, validatorService validator.Service) *UserHandler {
	return &UserHandler{dbService, validatorService}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	request := new(Dtos.CreateUserDto)
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(map[string]string{"message": "Invalid request payload"})
	}
	validationResult := services.Validate(h.validatorService.Validator, *h.validatorService.Trans, request)
	if validationResult != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(validationResult)
	}
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{"message": "Server Error"})
	}

	user := Models.User{
		Name:     request.Name,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: hashedPassword,
	}

	err = h.dbService.Db.Create(&user).Error
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
	validationResult := services.Validate(h.validatorService.Validator, *h.validatorService.Trans, request)
	if validationResult != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(validationResult)
	}
	user := Models.User{}
	h.dbService.Db.Where("phone = ?", request.Phone).First(&user)
	if user.ID == 0 {
		return c.Status(http.StatusUnauthorized).JSON(map[string]string{"message": "Invalid credentials"})
	}

	if !checkPasswordHash(request.Password, user.Password) {
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
	authService := services.NewAuthService(c, *h.dbService.Db)
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
