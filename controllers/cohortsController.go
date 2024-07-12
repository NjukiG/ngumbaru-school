package controllers

import (
	"github/NjukiG/ngumbaru-school/initializers"
	"github/NjukiG/ngumbaru-school/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Add a new cohort
func AddACohort(c *gin.Context) {
	var body struct {
		Name string
		Year int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to  read request body...",
		})
		return
	}

	cohort := models.Cohort{
		Name: body.Name,
		Year: body.Year,
	}
	// Get an admin to create
	user, _ := c.Get("user")

	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to create a cohort / Not admin",
		})
		return
	}

	result := initializers.DB.Create(&cohort)

	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Failed to create cohort...",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Cohort": cohort,
	})
}

// Func to get ALL COHORTS
func GetAllCohorts(c *gin.Context) {
	var cohorts []models.Cohort

	result := initializers.DB.Preload("Students").Preload("Courses").Preload("Teachers").Find(&cohorts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to fetch cohorts...",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Cohorts": cohorts})
}

// Get one cohort by Id

func GetCohortByID(c *gin.Context) {
	cohortId := c.Param("id")
	var cohort models.Cohort

	result := initializers.DB.Preload("Students").Preload("Courses").Preload("Teachers").First(&cohort, cohortId)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to fetch the cohort...",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Cohort": cohort})

}

// Enroll a student to a cohort
func EnrollStudentToCohort(c *gin.Context) {
	var body struct {
		StudentID uint
		CohortID  uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	var cohort models.Cohort
	result := initializers.DB.First(&cohort, body.CohortID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cohort not found"})
		return
	}

	var student models.User
	result2 := initializers.DB.First(&student, body.StudentID)

	if result2.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	user, _ := c.Get("user")

	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to add a student / Not an admin",
		})
		return
	}

	initializers.DB.Model(&cohort).Association("Students").Append(&student)

	c.JSON(http.StatusOK, gin.H{"Cohort": cohort})
}

// Add a teacher for a cohort
func AddCohortTeacher(c *gin.Context) {
	var body struct {
		TeacherID uint
		CohortID  uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	var cohort models.Cohort
	result := initializers.DB.First(&cohort, body.CohortID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cohort not found"})
		return
	}

	var teacher models.User
	result2 := initializers.DB.First(&teacher, body.TeacherID)

	if result2.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	user, _ := c.Get("user")

	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to add a teacher / Not an admin",
		})
		return
	}

	initializers.DB.Model(&cohort).Association("Teachers").Append(&teacher)
	c.JSON(http.StatusOK, gin.H{"Cohort": cohort})
}


func AddCourseToCohort(c *gin.Context) {
    var body struct {
        CohortID uint
        CourseID uint
    }

    if err := c.Bind(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
        return
    }

    var cohort models.Cohort
    if result := initializers.DB.First(&cohort, body.CohortID); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cohort not found"})
        return
    }

    var course models.Course
    if result := initializers.DB.First(&course, body.CourseID); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
        return
    }

    initializers.DB.Model(&cohort).Association("Courses").Append(&course)

    c.JSON(http.StatusOK, gin.H{"message": "Course added to cohort"})
}
