package models

type LoginRequest struct {
	Username string `json:"username" validate:"required"` // bisa username atau email
	Password string `json:"password" validate:"required"`
}
