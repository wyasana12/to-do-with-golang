package controllers

import (
	"encoding/json"
	"net/http"
	"time"
	"to-do-list-go/config"
	"to-do-list-go/helper"
	"to-do-list-go/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var register models.Register

	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	defer r.Body.Close()

	if err := config.Validate.Struct(register); err != nil {
		helper.Response(w, 400, "Validation Error", helper.FormatValidationError(err))
		return
	}

	if register.Password != register.PasswordConfirm {
		helper.Response(w, 400, "Password Not Match", nil)
		return
	}

	passwordHash, err := helper.HashPassword(register.Password)
	if err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	expired := time.Now().Add(5 * time.Minute)

	user := models.User{
		Name:              register.Name,
		Username:          register.Username,
		Email:             register.Email,
		Password:          passwordHash,
		VerificationToken: helper.GenerateVerificationToken(32),
		EmailVerified:     false,
		ResetTokenExpiry:  &expired,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	if err := helper.SendVerificationEmail(user.Email, user.VerificationToken); err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Success Register User", nil)
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		helper.Response(w, 400, "Missing Token", nil)
		return
	}

	var user models.User
	if err := config.DB.Where("verification_token = ?", token).First(&user).Error; err != nil {
		helper.Response(w, 404, "Invalid Or Expired Token", nil)
		return
	}

	if user.ResetTokenExpiry != nil && time.Now().After(*user.ResetTokenExpiry) {
		helper.Response(w, 400, "Verification Token Expired", nil)
		return
	}

	if user.EmailVerified {
		helper.Response(w, 404, "Email Already Verified", nil)
		return
	}

	user.EmailVerified = true
	user.VerificationToken = ""
	user.ResetTokenExpiry = nil
	if err := config.DB.Save(&user).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Email Success Verified", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var login models.Login

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	defer r.Body.Close()

	if err := config.Validate.Struct(login); err != nil {
		helper.Response(w, 400, "Validation Error", helper.FormatValidationError(err))
		return
	}

	var user models.User
	if err := config.DB.First(&user, "email = ? OR username = ?", login.Email, login.Username).Error; err != nil {
		helper.Response(w, 404, "Wrong Email/Username or Password", nil)
		return
	}

	if err := helper.VerifyPassword(user.Password, login.Password); err != nil {
		helper.Response(w, 404, "Wrong Email/Username or Password", nil)
		return
	}

	if !user.EmailVerified {
		helper.Response(w, 401, "Email Not Verified", nil)
		return
	}

	token, err := helper.CreateToken(&user)
	if err != nil {
		helper.Response(w, 404, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Successfully Login", token)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var reset models.Reset

	if err := json.NewDecoder(r.Body).Decode(&reset); err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	defer r.Body.Close()

	if err := config.Validate.Struct(reset); err != nil {
		helper.Response(w, 400, "Validation Error", helper.FormatValidationError(err))
		return
	}

	var user models.User
	if err := config.DB.Where("email = ? AND username = ?", reset.Email, reset.Username).First(&user).Error; err != nil {
		helper.Response(w, 404, "Wrong Email And Username", nil)
		return
	}

	hashed, err := helper.HashPassword(reset.NewPassword)

	if err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	user.Password = hashed
	config.DB.Save(&user)

	helper.Response(w, 201, "Success Reset Password", nil)
}
