package partners

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

func (h *Handler) AllPartnersHandler(c *gin.Context) {
	userCompanyID, exist := c.Get("company_id")
	uid := c.GetInt("user_id")
	if !exist {
		slog.Error("all_partners - company_id missing", "user_id", uid)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you don't have access to this page",
		})
		return
	}
	allUsers, err := h.service.AllPartners(userCompanyID.(int))

	if err != nil {
		slog.Error("all_partners - failed", "error", err, "company_id", userCompanyID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.Info("all_partners - success", "company_id", userCompanyID, "count", len(allUsers.Partners))
	c.JSON(http.StatusOK, *allUsers)

}

func (h *Handler) PartnerInfoHandler(c *gin.Context) {
	partnerId := c.Param("PartnerID")
	uid := c.GetInt("user_id")

	partner, err := h.service.ReadPartner(partnerId)

	if err != nil {
		slog.Error("partner_info - failed", "error", err, "partner_id", partnerId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("partner_info - success", "partner_id", partnerId, "request_by", uid)
	c.JSON(http.StatusOK, *partner)

}

func (h *Handler) PartnerCreationHandler(c *gin.Context) {
	var req CreatePartnerRequest
	uid := c.GetInt("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("create_partner - validation failed", "error", err, "request_by", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}

	comapnyID := c.GetInt("company_id")
	err := h.service.CreatePartner(&req, comapnyID)

	if err != nil {
		slog.Error("create_partner - failed", "error", err, "partner_name", req.Name, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("create_partner - success", "partner_name", req.Name, "request_by", uid)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Partner created succesfully",
	})

}

func (h *Handler) PartnerUpdateHandler(c *gin.Context) {
	var req UpdatePartnerRequest
	uid := c.GetInt("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("update_partner - validation failed", "error", err, "request_by", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}
	err := h.service.UpdatePartner(&req, req.ID)

	if err != nil {
		slog.Error("update_partner - failed", "error", err, "partner_id", req.ID, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	slog.Info("update_partner - success", "partner_id", req.ID, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Partner updated succesfully",
	})
}

func (h *Handler) PartnerDeleteHandler(c *gin.Context) {
	partnerId := c.Param("PartnerID")
	uid := c.GetInt("user_id")

	user, err := h.service.ReadPartner(partnerId)

	if err != nil {
		slog.Error("delete_partner - read failed", "error", err, "partner_id", partnerId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = h.service.DeletePartner(user.ID)
	if err != nil {
		slog.Error("delete_partner - failed", "error", err, "partner_id", user.ID, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("delete_partner - success", "partner_id", user.ID, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Partner deleted succesfully",
	})
}
