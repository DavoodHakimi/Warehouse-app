package users

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

func (h *Handler) AllUsersHandler(c *gin.Context) {
	userCompanyID, exist := c.Get("company_id")
	uid := c.GetInt("user_id")
	if !exist {
		slog.Error("all_users - company_id missing", "user_id", uid)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you don't have access to this page",
		})
		return
	}
	allUsers, err := h.service.AllUsers(userCompanyID.(int))

	if err != nil {
		slog.Error("all_users - failed", "error", err, "company_id", userCompanyID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.Info("all_users - success", "company_id", userCompanyID)
	c.JSON(http.StatusOK, *allUsers)

}

func (h *Handler) UserInfoHandler(c *gin.Context) {
	userId := c.Param("userID")
	uid := c.GetInt("user_id")
	companyID := c.GetInt("company_id")

	user, err := h.service.ReadUser(userId, companyID)

	if err != nil {
		slog.Error("user_info - failed", "error", err, "target_user_id", userId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("user_info - success", "target_user_id", userId, "request_by", uid)
	c.JSON(http.StatusOK, *user)

}

func (h *Handler) UserCreationHandler(c *gin.Context) {
	var req CreateUserRequest
	uid := c.GetInt("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("create_user - validation failed", "error", err, "request_by", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}
	comapnyID := c.GetInt("company_id")
	err := h.service.CreateUser(&req, comapnyID)

	if err != nil {
		slog.Error("create_user - failed", "error", err, "username", req.UserName, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("create_user - success", "username", req.UserName, "user_type", req.UserTypeID, "request_by", uid)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created succesfully",
	})

}

func (h *Handler) UserUpdateHandler(c *gin.Context) {
	var req UpdateUserRequest
	uid := c.GetInt("user_id")
	companyID := c.GetInt("company_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("update_user - validation failed", "error", err, "request_by", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}
	err := h.service.UpdateUser(&req, uid, companyID)

	if err != nil {
		slog.Error("update_user - failed", "error", err, "target_user_id", req.ID, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	slog.Info("update_user - success", "target_user_id", req.ID, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated succesfully",
	})
}

func (h *Handler) UserDeleteHandler(c *gin.Context) {
	userId := c.Param("userID")
	uid := c.GetInt("user_id")
	companyID := c.GetInt("company_id")

	user, err := h.service.ReadUser(userId, companyID)

	if err != nil {
		slog.Error("delete_user - read failed", "error", err, "target_user_id", userId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = h.service.DeleteUser(user.ID, companyID)
	if err != nil {
		slog.Error("delete_user - failed", "error", err, "target_user_id", user.ID, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("delete_user - success", "target_user_id", user.ID, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted succesfully",
	})
}
