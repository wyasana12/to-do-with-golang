package trashcontroller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"to-do-list-go/config"
	"to-do-list-go/helper"
	"to-do-list-go/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Trash(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)
	var todo []models.Todo

	if err := config.DB.Unscoped().Where("user_id = ? AND deleted_at IS NOT NULL", user.ID).Find(&todo).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	if len(todo) == 0 {
		helper.Response(w, 404, "No Trashed Todo Found", nil)
		return
	}

	helper.Response(w, 200, "List Trashed Todo", todo)
}

func Restore(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var todo models.Todo
	if err := config.DB.Unscoped().Where("id = ? AND user_id = ?", id, user.ID).Find(&todo).Error; err != nil {
		helper.Response(w, 404, "Todo Not Found", nil)
		return
	}

	if todo.DeletedAt.Valid {
		todo.DeletedAt = gorm.DeletedAt{}
		config.DB.Unscoped().Save(&todo)
		helper.Response(w, 200, "Restored Todo", todo)
		return
	}

	helper.Response(w, 400, "Todo Is Not Deleted", nil)
}

func PermanentDelete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	idParams := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParams)

	var todo models.Todo
	if err := config.DB.Unscoped().Where("id = ? AND user_id = ?", id, user.ID).Delete(&todo).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "Todo Permanently Deleted", nil)
}

func BulkRestore(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	var payload models.BulkRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || len(payload.IDs) == 0 {
		helper.Response(w, 400, "Invalid Or Empty Id List", nil)
		return
	}

	defer r.Body.Close()

	if err := config.Validate.Struct(payload); err != nil {
		helper.Response(w, 400, "Validation Error", helper.FormatValidationError(err))
		return
	}

	var todo []models.Todo
	if err := config.DB.Unscoped().Where("user_id = ? AND id IN ?", user.ID, payload.IDs).Find(&todo).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	if len(todo) == 0 {
		helper.Response(w, 404, "Not Matching Trashed Todos Found", nil)
		return
	}

	for i := range todo {
		if todo[i].DeletedAt.Valid {
			todo[i].DeletedAt = gorm.DeletedAt{}
			config.DB.Unscoped().Save(&todo[i])
		}
	}

	helper.Response(w, 200, "Todos Restored Successfully", nil)
}

func BulkDelete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	var payload models.BulkRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helper.Response(w, 400, "Invalid Or Empty List", nil)
		return
	}

	defer r.Body.Close()

	if err := config.Validate.Struct(payload); err != nil {
		helper.Response(w, 400, "Validation Error", helper.FormatValidationError(err))
		return
	}

	var todo []models.Todo

	res := config.DB.Unscoped().Where("user_id = ? AND id IN ?", user.ID, payload.IDs).Delete(&todo)
	if res.Error != nil {
		helper.Response(w, 500, res.Error.Error(), nil)
		return
	}

	if res.RowsAffected == 0 {
		helper.Response(w, 404, "No Matching Trashes Todos Found", nil)
		return
	}

	helper.Response(w, 200, "Todos Permanently Delete", nil)
}
