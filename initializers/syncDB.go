package initializers

import (
	"github/NjukiG/ngumbaru-school/models"
	"log"
)

func SyncDatabase() {
	err := DB.AutoMigrate(&models.User{}, &models.Course{}, &models.Cohort{}, &models.Attendance{})

	if err != nil {
		log.Fatal("Failed to migrate model", err)
	}
}
