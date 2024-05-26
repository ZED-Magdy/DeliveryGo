package Dtos

type LoginDto struct {
	Phone    string `json:"phone" validate:"required,numeric"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}
