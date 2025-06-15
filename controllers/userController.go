package controllers

import (
	"encoding/json"
	"net/http"
	"to-do-list-go/config"
	"to-do-list-go/helper"
	"to-do-list-go/models"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	var dbUser models.User
	if err := config.DB.First(&dbUser, user.ID).Error; err != nil {
		helper.Response(w, 404, "User Not Found", nil)
		return
	}

	profile := models.Profile{
		ID:       dbUser.ID,
		Name:     dbUser.Name,
		Username: dbUser.Username,
		Email:    dbUser.Email,
	}

	helper.Response(w, 200, "My Profile", profile)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	var input models.Profile
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.Response(w, 400, "Invalid Input", nil)
		return
	}

	defer r.Body.Close()

	var dbUser models.User
	if err := config.DB.First(&dbUser, user.ID).Error; err != nil {
		helper.Response(w, 404, "User Not Found", nil)
		return
	}

	dbUser.Name = input.Name
	dbUser.Username = input.Username
	dbUser.Email = input.Email

	if err := config.DB.Save(&dbUser).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Success Update Profile", nil)
}
