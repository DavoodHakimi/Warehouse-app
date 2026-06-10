package partners

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

func (h *Handler) AllPartnersHandler(c *gin.Context) {
	userCompanyID, exist := c.Get("company_id")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you don't have access to this page",
		})
		return
	}
	allUsers, err := h.service.AllPartners(userCompanyID.(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *allUsers)

}

func (h *Handler) PartnerInfoHandler(c *gin.Context) {
	partnerId := c.Param("PartnerID")

	partner, err := h.service.ReadPartner(partnerId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *partner)

}

func (h *Handler) PartnerCreationHandler(c *gin.Context) {
	var req CreatePartnerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
	}
	comapnyID, _ := c.Get("company_id")
	err := h.service.CreatePartner(&req, comapnyID.(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Partner created succesfully",
	})

}

func (h *Handler) PartnerUpdateHandler(c *gin.Context) {
	var req UpdatePartnerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}
	err := h.service.UpdatePartner(&req, req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Partner updated succesfully",
	})
}

func (h *Handler) PartnerDeleteHandler(c *gin.Context) {
	partnerId := c.Param("PartnerID")

	user, err := h.service.ReadPartner(partnerId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = h.service.DeletePartner(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Partner deleted succesfully",
	})
}
