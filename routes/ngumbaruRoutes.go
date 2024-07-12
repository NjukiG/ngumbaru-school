package routes

import (
	"github/NjukiG/ngumbaru-school/controllers"
	"github/NjukiG/ngumbaru-school/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	publicRoutes := router.Group("/public")
	{
		// ROutes to signup user and Login user. In public
		publicRoutes.POST("/register", controllers.RegisterUser)
		publicRoutes.POST("/login", controllers.LoginUser)
	}

	// For all protected routes,  user has to be logged in  and have correct middleware access to access them
	protectedRoutes := router.Group("/protected")
	protectedRoutes.Use(middleware.RequireAuth)

	{
		// ROutes to logout and validate a logged in user. In private.
		protectedRoutes.GET("/validate", controllers.ValidateUser)
		protectedRoutes.POST("/logout", controllers.LogoutUser)
		// Courses routes
		protectedRoutes.POST("/courses", controllers.CreateCourse)
		protectedRoutes.GET("/courses", controllers.GetAllCourses)
		protectedRoutes.GET("/courses/:id", controllers.GetCourseByID)
		protectedRoutes.PUT("/courses/:id", controllers.EditCourseDetails)
		protectedRoutes.POST("/courses/enrollStudent", controllers.EnrollStudent)
		protectedRoutes.POST("/courses/addTeacher", controllers.AddCourseTeacher)
		protectedRoutes.DELETE("/courses/:id", controllers.DeleteCourse)
		protectedRoutes.POST("/courses/addCohort", controllers.AddCohortToCourse)

		// Cohort Routes
		protectedRoutes.POST("/cohorts", controllers.AddACohort)
		protectedRoutes.GET("/cohorts", controllers.GetAllCohorts)
		protectedRoutes.GET("/cohorts/:id", controllers.GetCohortByID)
		protectedRoutes.POST("/cohorts/enrollStudent", controllers.EnrollStudentToCohort)
		protectedRoutes.POST("/cohorts/addTeacher", controllers.AddCohortTeacher)
		protectedRoutes.POST("cohorts/addCourse", controllers.AddCourseToCohort)

		// Attendance routes
		protectedRoutes.POST("/attendance", middleware.AdminOnly(), controllers.MarkAttendance)
		protectedRoutes.GET("/attendance", controllers.GetAttendance)
	}
}
