package Dtos

type UserDto struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Email *string `json:"email"`
	Phone string  `json:"phone"`
	Token *string  `json:"token"`
}
