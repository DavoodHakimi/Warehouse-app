package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AllUsersHandler(c *gin.Context) {
	userRole, exist := c.Get("role")
	if !exist || userRole != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you don't have access to this page",
		})
		return
	}
	userCompanyID, _ := c.Get("company_id")
	allUsers, err := h.service.AllUsers(userCompanyID.(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *allUsers)

}

func (h *Handler) UserInfoHandler(c *gin.Context) {
	userId := c.Param("userID")

	user, err := h.service.ReadUser(userId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *user)

}

func (h *Handler) UserCreationHandler(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
	}

	err := h.service.CreateUser(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created succesfully",
	})

}

func (h *Handler) UserUpdateHandler(c *gin.Context) {

}

func (h *Handler) UserDeleteHandler(c *gin.Context) {

}
