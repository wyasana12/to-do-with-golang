package controllers

import (
	"encoding/json"
	"net/http"
	"to-do-list-go/config"
	"to-do-list-go/helper"
	"to-do-list-go/models"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Register(w http.ResponseWriter, r *http.Request) {
	var register models.Register

	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	defer r.Body.Close()

	if err := validate.Struct(register); err != nil {
		helper.Response(w, 400, "Validation Not Match: "+err.Error(), nil)
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

	user := models.User{
		Name:     register.Name,
		Username: register.Username,
		Email:    register.Email,
		Password: passwordHash,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Success Register User", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var login models.Login

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		helper.Response(w, 500, err.Error(), nil)
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

	token, err := helper.CreateToken(&user)
	if err != nil {
		helper.Response(w, 404, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Successfully Login", token)
}

func Logout(w http.ResponseWriter, r *http.Request) {

}
