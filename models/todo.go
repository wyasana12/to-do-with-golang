package models

import (
	"time"

	"gorm.io/gorm"
)

type TodoStatus string

const (
	StatusNotStarted TodoStatus = "not_started"
	StatusInProgress TodoStatus = "in_progress"
	StatusCompleted  TodoStatus = "completed"
)

func (t TodoStatus) IsValid() bool {
	switch t {
	case StatusNotStarted, StatusInProgress, StatusCompleted:
		return true
	}
	return false
}

type Todo struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	Title                 string         `json:"title" validate:"required,min=3,max=100"`
	UserID                uint           `json:"user_id" validate:"required"`
	User                  User           `gorm:"foreignKey:UserID" json:"user"`
	Description           string         `json:"description" validate:"omitempty,max=1000"`
	Status                TodoStatus     `gorm:"type:varchar(50);default:'not_started'" json:"status" validate:"omitempty,oneof=not_started in_progress completed"`
	StartDate             *time.Time     `json:"start_date" validate:"omitempty,ltefield=EndDate"`
	EndDate               *time.Time     `json:"end_date" validate:"omitempty,gtefield=StartDate"`
	IsD1Notified          bool           `gorm:"default:false" json:"is_d1_notified"`
	IsLessThan1HrNotified bool           `gorm:"default:false" json:"is_less_than_1hr_notified"`
	IsOverdueDeadline     bool           `gorm:"default:false" json:"is_overdue_notified"`
	Attachments           []Attachment   `gorm:"foreignKey:TodoID" json:"attachments"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}

type TodoResponse struct {
	ID                    uint                 `json:"id"`
	Title                 string               `json:"title"`
	UserID                uint                 `json:"-"`
	User                  Profile              `json:"user"`
	Description           string               `json:"description"`
	Status                TodoStatus           `json:"status"`
	StartDate             time.Time            `json:"start_date"`
	EndDate               time.Time            `json:"end_date"`
	IsD1Notified          bool                 `json:"is_d1_notified"`
	IsLessThan1HrNotified bool                 `json:"is_less_than_1hr_notified"`
	IsOverdueDeadline     bool                 `json:"is_overdue_notified"`
	Attachments           []AttachmentResponse `json:"attachments"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
}

type UpdateTodoRequest struct {
	Title       string     `json:"title" validate:"required,min=3,max=100"`
	Description string     `json:"description" validate:"omitempty,max=1000"`
	Status      TodoStatus `gorm:"type:varchar(50);default:'not_started'" json:"status" validate:"omitempty,oneof=not_started in_progress completed"`
	StartDate   *time.Time `json:"start_date" validate:"omitempty,ltefield=EndDate"`
	EndDate     *time.Time `json:"end_date" validate:"omitempty,gtefield=StartDate"`
}

type BulkRequest struct {
	IDs []uint `json:"ids" validate:"required,min=1,dive,gt=0"`
}
