package controllers

import (
	"github/NjukiG/ngumbaru-school/initializers"
	"github/NjukiG/ngumbaru-school/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Func to create a course. Only by an admin
func CreateCourse(c *gin.Context) {
	var body struct {
		Name     string
		SubTitle string
		ImageURL string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read request body",
		})
		return
	}

	course := models.Course{
		Name:     body.Name,
		SubTitle: body.SubTitle,
		ImageURL: body.ImageURL,
	}

	// Get an admin to create
	user, _ := c.Get("user")

	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to create a course / Not admin",
		})
		return
	}

	result := initializers.DB.Create(&course)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create course",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Course": course,
	})
}

// Func to get all courses.
func GetAllCourses(c *gin.Context) {
	var courses []models.Course

	result := initializers.DB.Preload("Teachers").Preload("Students").Find(&courses)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to fetch courses",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Courses": courses})
}

// Func to get a single course by ID
func GetCourseByID(c *gin.Context) {
	courseId := c.Param("id")
	var course models.Course

	result := initializers.DB.Preload("Teachers").Preload("Students").First(&course, courseId)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Course not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Course": course})
}

// Func to edit a course details by admin only
func EditCourseDetails(c *gin.Context) {
	courseId := c.Param("id")

	var body struct {
		Name     string
		SubTitle string
		ImageURL string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	var course models.Course

	initializers.DB.First(&course, courseId)

	user, _ := c.Get("user")

	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to edit course details",
		})
		return
	}

	initializers.DB.Model(&course).Updates(models.Course{
		Name:     body.Name,
		SubTitle: body.SubTitle,
		ImageURL: body.ImageURL,
	})

	c.JSON(http.StatusOK, course)
}

// FUnc to delete a course by Admin only

func DeleteCourse(c *gin.Context) {

	courseId := c.Param("id")

	var course models.Course

	result := initializers.DB.First(&course, courseId)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Course not found",
		})
		return
	}

	user, _ := c.Get("user")

	if user.(models.User).Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"Error": "Not allowed to delete a course / Not an admin",
		})
		return
	}

	initializers.DB.Delete(&course, courseId)

	// Respond
	c.Status(http.StatusNoContent)

	c.JSON(200, gin.H{
		"Message": "A course was deleted...",
	})

}

// Enroll a student in a course
func EnrollStudent(c *gin.Context) {
	var body struct {
		StudentID uint
		CourseID  uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	var course models.Course
	result := initializers.DB.First(&course, body.CourseID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
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

	initializers.DB.Model(&course).Association("Students").Append(&student)

	c.JSON(http.StatusOK, course)
}

// Add a Teacher for the course
func AddCourseTeacher(c *gin.Context) {
	var body struct {
		TeacherID uint
		CourseID  uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	var course models.Course
	result := initializers.DB.First(&course, body.CourseID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
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

	initializers.DB.Model(&course).Association("Teachers").Append(&teacher)

	c.JSON(http.StatusOK, course)
}


func AddCohortToCourse(c *gin.Context) {
    var body struct {
        CourseID uint
        CohortID uint
    }

    if err := c.Bind(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
        return
    }

    var course models.Course
    if result := initializers.DB.First(&course, body.CourseID); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
        return
    }

    var cohort models.Cohort
    if result := initializers.DB.First(&cohort, body.CohortID); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cohort not found"})
        return
    }

    initializers.DB.Model(&course).Association("Cohorts").Append(&cohort)

    c.JSON(http.StatusOK, gin.H{"message": "Cohort added to course"})
}