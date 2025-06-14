package todocontroller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"to-do-list-go/config"
	"to-do-list-go/helper"
	"to-do-list-go/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Index(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	var todo []models.Todo

	if err := config.DB.Where("user_id = ?", user.ID).Find(&todo).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "List To-Do", &todo)
}

func Create(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)
	var todo models.Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	defer r.Body.Close()

	todo.UserID = user.ID
	todo.Completed = false

	if err := config.DB.Create(&todo).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 201, "Success Create ToDo", nil)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var todo models.Todo

	if err := config.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "ID ToDo Not Found or Invalid User Id", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "Detail ToDo", &todo)
}

func Update(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var todo models.Todo

	if err := config.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "ID ToDo Not Found or Invalid User Id", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	var updatedData models.Todo

	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		helper.Response(w, 400, "Invalid JSON Body", nil)
		return
	}

	defer r.Body.Close()

	todo.Title = updatedData.Title
	todo.Completed = updatedData.Completed

	if err := config.DB.Save(&todo).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "Success Updated ToDo", nil)
}

func Destroy(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var todo models.Todo

	if err := config.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "ID ToDo Not Found or Invalid User Id", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	if err := config.DB.Delete(&todo).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "Success Delete ToDo", nil)
}
