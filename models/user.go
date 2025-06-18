package models

import "time"

type User struct {
	ID                 uint `gorm:"primaryKey"`
	Name               string
	Username           string `gorm:"unique"`
	Email              string `gorm:"unique"`
	Password           string
	EmailVerified      bool       `gorm:"default:false"`
	VerificationToken  string     `gorm:"index"`
	ResetPasswordToken string     `gorm:"index"`
	ResetTokenExpiry   *time.Time `gorm:"default:null"`
	// OTPCode            string     `gorm:"type:varchar(6);index"`
	// OTPExpiry          *time.Time `gorm:"default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Register struct {
	Name            string `json:"name" validate:"required"`
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}

type Login struct {
	Username string `json:"username" validate:"omitempty,required_without=Email"`
	Email    string `json:"email" validate:"omitempty,required_without=Username,email"`
	Password string `json:"password" validate:"required"`
}

type Profile struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Reset struct {
	Email       string `json:"email" validate:"required,email"`
	Username    string `json:"username" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}
