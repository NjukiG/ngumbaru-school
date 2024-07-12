package controllers

import (
	"github/NjukiG/ngumbaru-school/initializers"
	"github/NjukiG/ngumbaru-school/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	//  Get Nname, email & password of req user
	var body struct {
		FirstName string
		LastName  string
		Role      models.Role
		Email     string
		Password  string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read req body...",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Hash the password
	// Create the user
	user := models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Role:      body.Role,
		Email:     body.Email,
		Password:  string(hash),
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create a new user. Use a different email Address",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "New user created",
	})
}

// Function to Login a user.
func LoginUser(c *gin.Context) {
	// Get the email and password of req user
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read req body...",
		})
		return
	}

	// Lookup requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})
		return
	}
	// Compare send in password with hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})
		return
	}

	// Generate a JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	// Send it back
	c.JSON(http.StatusOK, gin.H{
		"Logged In?": true,
		// "token":      tokenString,   // Dont Show token
	})
}

// Function to  validate that user is actually logged in.
func ValidateUser(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"Message": "I am Logged in...",
		"User":    user,
	})
}

// Logout function
func LogoutUser(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
