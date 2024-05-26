package Handlers

import (
	"ZED-Magdy/Delivery-go/Dtos"
	"ZED-Magdy/Delivery-go/Models"
	"ZED-Magdy/Delivery-go/Services"
	"encoding/json"
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	request := Dtos.CreateUserDto{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}
	err = h.validator.Struct(request)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Error validating request"}`))
			return
		}
		errs := err.(validator.ValidationErrors)
		errors := errs.Translate(h.trans)

		w.WriteHeader(http.StatusUnprocessableEntity)
		errors_marshalled, _ := json.Marshal(errors)
		w.Write(errors_marshalled)
		return

	}
	hashedPassword, err := h.hashPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
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
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "User already exists"}`))
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Error creating user"}`))
				return

			}
		}
	}
	token, err := services.CreateJwtToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error creating token"}`))
		return
	}
	userDto := Dtos.UserDto{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Token: &token,
	}
	user_marshalled, _ := json.Marshal(userDto)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(user_marshalled))

}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	request := Dtos.LoginDto{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload"}`))
		return
	}
	err = h.validator.Struct(request)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Error validating request"}`))
			return
		}
		errs := err.(validator.ValidationErrors)
		errors := errs.Translate(h.trans)

		w.WriteHeader(http.StatusUnprocessableEntity)
		errors_marshalled, _ := json.Marshal(errors)
		w.Write(errors_marshalled)
		return

	}
	user := Models.User{}
	h.db.Where("phone = ?", request.Phone).First(&user)
	if user.ID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid credentials"}`))
		return
	}

	if !h.checkPasswordHash(request.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid credentials"}`))
		return
	}

	token, err := services.CreateJwtToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error creating token"}`))
		return
	}
	userDto := Dtos.UserDto{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Token: &token,
	}
	user_marshalled, _ := json.Marshal(userDto)
	w.Write([]byte(user_marshalled))
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Unauthorized"}`))
		return
	}

	token := authHeader[1]
	tokenObj, err := services.VerifyJwtToken(token)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Unauthorized"}`))
		return
	}

	user := Models.User{}
	err = h.db.First(&user, tokenObj.Claims.(jwt.MapClaims)["subject"].(float64)).Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	userDto := Dtos.UserDto{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}
	user_marshalled, _ := json.Marshal(userDto)

	w.Write([]byte(user_marshalled))
}


func (*UserHandler) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return string(bytes), err
}

func (*UserHandler) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
