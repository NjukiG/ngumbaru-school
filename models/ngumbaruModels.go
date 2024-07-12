package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin   Role = "Admin"
	RoleTeacher Role = "Teacher"
	RoleStudent Role = "Student"
)

// Users Model: Admin, Teachers and Students
type User struct {
	gorm.Model
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Role      Role   `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `json:"-"`
}

// Courses Model
type Course struct {
	gorm.Model
	Name     string `gorm:"not null"`
	SubTitle string
	ImageURL string
	Teachers []*User   `gorm:"many2many:course_teachers"`
	Students []*User   `gorm:"many2many:course_students"`
	Cohorts  []*Cohort `gorm:"many2many:course_cohorts"`
}

// Cohorts model
// The name can be a code t identify the cohort e.g SD 59-63
type Cohort struct {
	gorm.Model
	Name     string    `gorm:"not null"`
	Year     int       `gorm:"not null"`
	Students []*User   `gorm:"many2many:cohort_students"`
	Courses  []*Course `gorm:"many2many:cohort_courses"`
	Teachers []*User   `gorm:"many2many:cohort_teachers"`
}

// Attendance model to check student attendance
type Attendance struct {
	gorm.Model
	Date      time.Time `gorm:"not null"`
	StudentID uint      `gorm:"not null"`
	CohortID  uint      `gorm:"not null"`
	Present   bool      `gorm:"not null"`
}
