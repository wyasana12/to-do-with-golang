package todocontroller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
	"to-do-list-go/config"
	"to-do-list-go/helper"
	"to-do-list-go/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func derefTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}

	return time.Time{}
}

func Index(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)
	query := config.DB.Model(&models.Todo{}).Preload("User").Where("user_id = ?", user.ID)

	//filtering
	if status := r.URL.Query().Get("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if title := r.URL.Query().Get("title"); title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	//sorting
	sortBy := r.URL.Query().Get("sort_by")

	if sortBy == "" {
		sortBy = "created_at"
	}

	order := r.URL.Query().Get("order")

	if order != "asc" {
		order = "desc"
	}

	query = query.Order(sortBy + " " + order)

	//pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	var todo []models.Todo

	if err := query.Offset(offset).Limit(limit).Find(&todo).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	var todoResponse []models.TodoResponse
	for _, todo := range todo {
		todoResponse = append(todoResponse, models.TodoResponse{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status,
			StartDate:   derefTime(todo.StartDate),
			EndDate:     derefTime(todo.EndDate),
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
			User: models.Profile{
				ID:       todo.User.ID,
				Name:     todo.User.Name,
				Username: todo.User.Username,
				Email:    todo.User.Email,
			},
		})
	}

	helper.Response(w, 200, "List To-Do", &todoResponse)
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

	if err := config.Validate.Struct(todo); err != nil {
		helper.Response(w, 400, "Validation Error", helper.FormatValidationError(err))
		return
	}

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

	if err := config.DB.Preload("User").Where("id = ? AND user_id = ?", id, user.ID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "ID Not Found or Invalid User Id", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	todoResponse := models.TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		User:        models.Profile{ID: todo.User.ID, Name: todo.User.Name},
		Status:      todo.Status,
		StartDate:   derefTime(todo.StartDate),
		EndDate:     derefTime(todo.EndDate),
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}

	helper.Response(w, 200, "Detail ToDo", todoResponse)
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

	var updatedData models.UpdateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		helper.Response(w, 400, "Invalid JSON Body", nil)
		return
	}

	defer r.Body.Close()

	if err := config.Validate.Struct(updatedData); err != nil {
		helper.Response(w, 400, "Validation Error", helper.FormatValidationError(err))
		return
	}

	oldEndDate := todo.EndDate
	oldStatus := todo.Status

	todo.Title = updatedData.Title
	todo.Description = updatedData.Description
	todo.Status = updatedData.Status
	todo.StartDate = updatedData.StartDate
	todo.EndDate = updatedData.EndDate

	if oldEndDate != nil && updatedData.EndDate != nil && !updatedData.EndDate.Equal(*oldEndDate) {
		if now := time.Now(); updatedData.EndDate.After(now) {
			todo.IsD1Notified = false
			todo.IsLessThan1HrNotified = false
			todo.IsOverdueDeadline = false
		}
	}

	if oldStatus == models.StatusCompleted && (updatedData.Status == models.StatusNotStarted || updatedData.Status == models.StatusInProgress) {
		todo.IsD1Notified = false
		todo.IsLessThan1HrNotified = false
		todo.IsOverdueDeadline = false
	}

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

func BulkDestroy(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	var payload models.BulkRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || len(payload.IDs) == 0 {
		helper.Response(w, 400, "Invalid or Empty Id List", nil)
		return
	}

	defer r.Body.Close()

	if err := config.Validate.Struct(payload); err != nil {
		helper.Response(w, 400, "Validation Error", helper.FormatValidationError(err))
		return
	}

	var todo models.Todo
	if err := config.DB.Where("user_id = ? AND id IN ?", user.ID, payload.IDs).Delete(&todo).Error; err != nil {
		helper.Response(w, 500, err.Error(), nil)
		return
	}

	helper.Response(w, 200, "Bulk Delete Success", nil)
}
