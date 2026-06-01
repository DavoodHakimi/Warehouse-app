package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {

	var req SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
	}

	if err := SignUp(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})

}

func LogInHandler(c *gin.Context) {

	var req LogInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
	}

	token, err, status := Login(&req)
	if err != nil && status == 404 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	if err != nil && status == 401 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "wrong password"})
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"token": token,
		})

}

func MeHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	companyID, ok := c.Get("company_id")
	userName, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "key's does not exist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"username":   userName,
		"company_id": companyID,
	})
}
