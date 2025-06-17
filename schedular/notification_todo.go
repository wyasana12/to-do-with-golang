package schedular

import (
	"time"
	"to-do-list-go/config"
	"to-do-list-go/helper"
	"to-do-list-go/models"

	log "github.com/sirupsen/logrus"
)

func StartNotificationSchedular() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			log.Info("Running todo notification check...")

			var todo []models.Todo

			if err := config.DB.Preload("User").Where("status != ? AND end_date IS NOT NULL", models.StatusCompleted).Find(&todo).Error; err != nil {
				log.Errorf("Error Fetching todos for notification check: %v", err)
				continue
			}

			now := time.Now()

			for _, todo := range todo {
				if todo.EndDate == nil || todo.EndDate.IsZero() {
					continue
				}

				userEmail := todo.User.Email

				if userEmail == "" {
					log.Warnf("Skipping notification for todo ID %d: User email not found", todo.ID)
					continue
				}

				timeLeft := todo.EndDate.Sub(now)

				//condition 1: 24 hours before deadline
				if timeLeft <= 24*time.Hour && timeLeft > 1*time.Hour && !todo.IsD1Notified {
					log.Infof("Sending d-1 notification for todo ID: %d, Title: %s", todo.ID, todo.Title)
					if err := helper.SendTodoReminderEmail(userEmail, &todo, "D-1"); err != nil {
						log.Errorf("Failed to send d-1 notification for todo ID %d: %v", todo.ID, err)
					} else {
						todo.IsD1Notified = true
						if err := config.DB.Save(&todo).Error; err != nil {
							log.Errorf("Failed to update IsD1Notified for todo ID %d: %v", todo.ID, err)
						}
					}
				}

				//condition 2: less than one hour before deadline
				if timeLeft <= 1*time.Hour && timeLeft > 0 && !todo.IsLessThan1HrNotified {
					log.Infof("Sending less than one hour notification for todo ID: %d, Title: %s", todo.ID, todo.Title)
					if err := helper.SendTodoReminderEmail(userEmail, &todo, "LessThan1Hr"); err != nil {
						log.Errorf("Failed to send less than one hour notification for todo ID %d: %v", todo.ID, err)
					} else {
						todo.IsLessThan1HrNotified = true
						if err := config.DB.Save(&todo).Error; err != nil {
							log.Errorf("Failed to update IsLessThan1HrNotified for todo ID %d: %v", todo.ID, err)
						}
					}
				}

				//condition 3: overdue and not completed
				if now.After(*todo.EndDate) && todo.Status != models.StatusCompleted && !todo.IsOverdueDeadline {
					log.Infof("Sending overdue notification for todo ID: %d, Title: %s", todo.ID, todo.Title)
					if err := helper.SendTodoReminderEmail(userEmail, &todo, "Overdue"); err != nil {
						log.Errorf("Failed to send overdue notification for todo ID %d: %v", todo.ID, err)
					} else {
						todo.IsOverdueDeadline = true
						if err := config.DB.Save(&todo).Error; err != nil {
							log.Errorf("Failed to update IsOverdueDeadline for todo ID %d: %v", todo.ID, err)
						}
					}
				}
			}

		}
	}()
}
