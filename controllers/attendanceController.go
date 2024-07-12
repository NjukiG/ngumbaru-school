package controllers

import (
	"github/NjukiG/ngumbaru-school/initializers"
	"github/NjukiG/ngumbaru-school/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MarkAttendance(c *gin.Context) {

	var body struct {
		StudentID uint
		CohortID  uint
		Present   bool
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to red request body...",
		})
		return
	}

	attendance := models.Attendance{
		Date:      time.Now(),
		StudentID: body.StudentID,
		CohortID:  body.CohortID,
		Present:   body.Present,
	}

	result := initializers.DB.Create(&attendance)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to mark attendance",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"attendance": attendance})
}

func GetAttendance(c *gin.Context) {
	var attendances []models.Attendance
	initializers.DB.Preload("Student").Preload("Cohort").Find(&attendances)

	c.JSON(http.StatusOK, gin.H{"attendances": attendances})
}
