package Dtos

type CreateUserDto struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    *string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,numeric"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}
