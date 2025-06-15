package models

import "time"

type User struct {
	ID                 int `gorm:"primaryKey"`
	Name               string
	Username           string `gorm:"unique"`
	Email              string `gorm:"unique"`
	Password           string
	EmailVerified      bool       `gorm:"default:false"`
	VerificationToken  string     `gorm:"index"`
	ResetPasswordToken string     `gorm:"index"`
	ResetTokenExpiry   *time.Time `gorm:"default:null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Register struct {
	Name            string `json:"name" validate:"required"`
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}

type Login struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Profile struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Reset struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	NewPassword string `json:"new_password"`
}
