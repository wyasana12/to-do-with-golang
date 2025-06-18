package attachmentcontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"to-do-list-go/config"
	"to-do-list-go/helper"
	"to-do-list-go/models"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const maxUploadSize = 10 << 20

func Upload(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	idParams := mux.Vars(r)["id"]
	todoID, err := strconv.Atoi(idParams)
	if err != nil {
		helper.Response(w, 400, "Invalid ToDo Id", err)
		return
	}

	var todo models.Todo
	if err := config.DB.Where("id = ? AND user_id = ?", todoID, user.ID).First(&todo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "ToDo Not Found Or Unauthorized", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	var attachment models.Attachment
	attachment.TodoID = uint(todoID)

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			helper.Response(w, 400, "File Too Large or Invalid Form Data", nil)
			return
		}

		fileURL, fileName, mimeType, err := helper.UploadFile(r, "file", "./uploads")
		if err != nil {
			if err == http.ErrMissingFile {
				helper.Response(w, 400, "No File Provided In Multipart Form Or Invalid Field Name 'file'", nil)
			} else {
				log.Errorf("Error Uploading File: %v", err)
				helper.Response(w, 500, "Failed To Upload File: "+err.Error(), nil)
			}
			return
		}

		if fileURL == "" {
			helper.Response(w, 400, "No File Provided In Multipart Form", nil)
			return
		}

		attachment.URL = fileURL
		attachment.FileName = fileName
		attachment.MimeType = mimeType

		if strings.HasPrefix(mimeType, "image/") {
			attachment.Type = models.AttachmentTypeImage
		} else {
			attachment.Type = models.AttachmentTypeFile
		}
	} else if strings.HasPrefix(contentType, "application/json") {
		var linkReq models.CreateAttachmentRequest

		if err := json.NewDecoder(r.Body).Decode(&linkReq); err != nil {
			helper.Response(w, 400, "Invalid JSON Body For Link Attachment", nil)
			return
		}

		defer r.Body.Close()

		if err := config.Validate.Struct(linkReq); err != nil {
			helper.Response(w, 400, "Validation Error: ", helper.FormatValidationError(err))
			return
		}

		if linkReq.Type != models.AttachmentTypeLink {
			helper.Response(w, 400, "JSON Attachment Type Must Be A 'link'", nil)
			return
		}

		attachment.URL = linkReq.URL
		// attachment.FileName = linkReq.FileName
		// attachment.MimeType = linkReq.MimeType
		attachment.Type = models.AttachmentTypeLink
	} else {
		helper.Response(w, 400, "Unsupported Content-Type. Use Multipart/form-data For Files Or application/json for Links", nil)
		return
	}

	if err := config.DB.Create(&attachment).Error; err != nil {
		log.Errorf("Error Saving Attachment to DB: %v", err)
		helper.Response(w, 500, "Failed To Save Attachment Details", nil)
		return
	}

	helper.Response(w, 201, "Attachment Uploaded/Added Successfully", attachment)
}

func Download(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	idParams := mux.Vars(r)["attachment_id"]
	id, err := strconv.Atoi(idParams)
	if err != nil {
		helper.Response(w, 400, "Invalid Attachment ID", nil)
		return
	}

	var attachment models.Attachment
	if err := config.DB.Preload("Todo").Where("id = ?", id).First(&attachment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Attachment Not Found", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	if attachment.Todo.UserID != user.ID {
		helper.Response(w, 401, "Unauthorized To Access This Attachment", nil)
		return
	}

	if attachment.Type != models.AttachmentTypeFile && attachment.Type != models.AttachmentTypeImage {
		helper.Response(w, 400, "Only Files And Images Can Be Download", nil)
		return
	}

	localPath := strings.TrimPrefix(attachment.URL, "/uploads/")
	fullPath := filepath.Join("./uploads", localPath)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Errorf("File Not Found On Disk For Attachment ID %d: %v", attachment.ID, err)
			helper.Response(w, 404, "File Not Found On Server", nil)
		} else {
			log.Errorf("Error Opening File For Attachment ID %d: %v", attachment.ID, err)
			helper.Response(w, 500, err.Error(), nil)
		}
	}

	defer file.Close()

	contentType := attachment.MimeType
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(attachment.FileName))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachment.FileName))

	if _, err := io.Copy(w, file); err != nil {
		log.Errorf("Error Streaming File Content For Attachment ID %d: %v", attachment.ID, err)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User Info").(*helper.MyCustomClass)

	idParams := mux.Vars(r)
	todoID, err := strconv.Atoi(idParams["id"])
	if err != nil {
		helper.Response(w, 404, "Invalid ToDo ID", nil)
		return
	}

	attachmentID, err := strconv.Atoi(idParams["attachment_id"])
	if err != nil {
		helper.Response(w, 404, "Invalid Attachment ID", nil)
		return
	}

	var attachment models.Attachment
	if err := config.DB.Preload("Todo").Where("id = ? AND todo_id = ?", attachmentID, todoID).First(&attachment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.Response(w, 404, "Attachment Not Found or Unautorizhed", nil)
			return
		}

		helper.Response(w, 500, err.Error(), nil)
		return
	}

	if attachment.Todo.UserID != user.ID {
		helper.Response(w, 401, "Unautorizhed To Delete This Attachment", nil)
		return
	}

	if attachment.Type == models.AttachmentTypeFile || attachment.Type == models.AttachmentTypeImage {
		localPath := strings.TrimPrefix(attachment.URL, "/uploads/")
		fullPath := filepath.Join("./uploads", localPath)

		if _, err := os.Stat(fullPath); err == nil {
			if err := os.Remove(fullPath); err != nil {
				log.Errorf("Failed To Delete File From Storage %s: %v", fullPath, err)
			} else {
				log.Infof("File %s Successfully Delete From Storage", fullPath)
			}
		} else if os.IsNotExist(err) {
			log.Warnf("File %s Not Found On Disk, But DB Record Exists. Proceeding With DB Deletion", fullPath)
		} else {
			log.Errorf("Error Checking File Existence %s: %v", fullPath, err)
		}
	}

	if err := config.DB.Delete(&attachment).Error; err != nil {
		log.Errorf("Failed To Delete Attachment From DB: %v", err)
		helper.Response(w, 500, "Failed To Delete Attachment Record", nil)
		return
	}

	helper.Response(w, 200, "Success Delete Attachment", nil)
}
