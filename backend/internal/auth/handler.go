package auth

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignUpHandler(c *gin.Context) {

	var req SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("signup - validation failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.SignUp(&req); err != nil {
		slog.Error("signup - failed", "error", err, "company", req.CompanyName)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("signup - success", "company", req.CompanyName, "user", req.UserName)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})

}

func (h *Handler) LogInHandler(c *gin.Context) {

	var req LogInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("login - validation failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}

	token, err, status := h.service.Login(&req)
	if err != nil && status == 404 {
		slog.Warn("login - user not found", "username", req.UserName)
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	if err != nil && status == 401 {
		slog.Warn("login - wrong password", "username", req.UserName)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "wrong password"})
		return
	}

	slog.Info("login - success", "username", req.UserName)
	c.JSON(
		http.StatusOK, gin.H{
			"token": token,
		})

}

func MeHandler(c *gin.Context) {
	userID, ok := c.Get("user_id")
	companyID, ok2 := c.Get("company_id")
	userName, ok3 := c.Get("username")
	if !ok || !ok2 || !ok3 {
		slog.Error("me - missing context keys")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "key's does not exist",
		})
		return
	}
	uid := userID.(int)
	slog.Info("me - success", "user_id", uid)
	c.JSON(http.StatusOK, meResponse{
		UserID:    uid,
		UserName:  userName.(string),
		CompanyID: companyID.(int),
	})
}
