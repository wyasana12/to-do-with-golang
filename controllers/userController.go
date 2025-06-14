package controllers

import (
	"net/http"
	"to-do-list-go/helper"
	"to-do-list-go/models"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)
	userResponse := &models.Profile{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}

	helper.Response(w, 200, "My Profile", userResponse)
}
