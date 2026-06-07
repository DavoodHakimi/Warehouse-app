package products

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

func (h *Handler) AllProductsHandler(c *gin.Context) {
	userRole, exist := c.Get("role")
	if !exist || userRole != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you don't have access to this page",
		})
		return
	}
	userCompanyID, _ := c.Get("company_id")
	allProducts, err := h.service.AllProducts(userCompanyID.(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *allProducts)

}

func (h *Handler) ProductInfoHandler(c *gin.Context) {
	productNum := c.Param("productNumber")

	prod, err := h.service.ReadProduct(productNum)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *prod)

}

func (h *Handler) ProductCreationHandler(c *gin.Context) {
	var req ProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
	}
	comapnyID, _ := c.Get("company_id")
	err := h.service.CreateProduct(&req, comapnyID.(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created succesfully",
	})

}

func (h *Handler) ProductUpdateHandler(c *gin.Context) {
	var req UpdateProductRequest

	userID, exist := c.Get("user_id")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}
	err := h.service.UpdateProduct(&req, userID.(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated succesfully",
	})
}

func (h *Handler) ProductDeleteHandler(c *gin.Context) {
	productId := c.Param("productNumber")

	prod, err := h.service.ReadProduct(productId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = h.service.DeleteProduct(prod.ProductNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted succesfully",
	})
}
