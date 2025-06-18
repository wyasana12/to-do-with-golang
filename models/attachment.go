package models

import "time"

type AttachmentType string

const (
	AttachmentTypeImage AttachmentType = "image"
	AttachmentTypeFile  AttachmentType = "file"
	AttachmentTypeLink  AttachmentType = "link"
)

type Attachment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TodoID    uint           `json:"todo_id"`
	Todo      Todo           `gorm:"foreignKey:TodoID"`
	FileName  string         `json:"file_name"`
	URL       string         `json:"url"`
	MimeType  string         `json:"mime_type"`
	Type      AttachmentType `gorm:"type:varchar(50)" json:"type"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type CreateAttachmentRequest struct {
	// TodoID   uint           `json:"todo_id" validate:"required"`
	FileName string         `json:"file_name" validate:"omitempty"`
	URL      string         `json:"url" validate:"required,url"`
	MimeType string         `json:"mime_type" validate:"omitempty"`
	Type     AttachmentType `json:"type" validate:"required,oneof=image file link"`
}

type AttachmentResponse struct {
	ID        uint           `json:"id"`
	TodoID    uint           `json:"-"`
	Todo      Todo           `json:"-"`
	FileName  string         `json:"file_name"`
	URL       string         `json:"url"`
	MimeType  string         `json:"mime_type"`
	Type      AttachmentType `json:"type"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
