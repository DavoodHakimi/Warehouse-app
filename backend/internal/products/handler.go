package products

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

func (h *Handler) AllProductsHandler(c *gin.Context) {
	userCompanyID, exist := c.Get("company_id")
	uid := c.GetInt("user_id")
	if !exist {
		slog.Error("all_products - company_id missing", "user_id", uid)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you don't have access to this page",
		})
		return
	}
	allProducts, err := h.service.AllProducts(userCompanyID.(int))

	if err != nil {
		slog.Error("all_products - failed", "error", err, "company_id", userCompanyID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.Info("all_products - success", "company_id", userCompanyID, "count", len(allProducts.Products))
	c.JSON(http.StatusOK, *allProducts)

}

func (h *Handler) ProductInfoHandler(c *gin.Context) {
	productNum := c.Param("productNumber")
	uid := c.GetInt("user_id")
	companyID := c.GetInt("company_id")

	prod, err := h.service.ReadProduct(productNum, companyID)

	if err != nil {
		slog.Error("product_info - failed", "error", err, "product_number", productNum, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("product_info - success", "product_number", productNum, "request_by", uid)
	c.JSON(http.StatusOK, *prod)

}

func (h *Handler) ProductCreationHandler(c *gin.Context) {
	var req ProductRequest
	uid := c.GetInt("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("create_product - validation failed", "error", err, "request_by", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}
	comapnyID := c.GetInt("company_id")
	err := h.service.CreateProduct(&req, comapnyID)

	if err != nil {
		slog.Error("create_product - failed", "error", err, "product_name", req.Name, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("create_product - success", "product_name", req.Name, "request_by", uid)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created succesfully",
	})

}

func (h *Handler) ProductUpdateHandler(c *gin.Context) {
	var req UpdateProductRequest
	uid := c.GetInt("user_id")
	companyID := c.GetInt("company_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("update_product - validation failed", "error", err, "request_by", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}
	err := h.service.UpdateProduct(&req, uid, companyID)

	if err != nil {
		slog.Error("update_product - failed", "error", err, "product_number", req.ProductNumber, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	slog.Info("update_product - success", "product_number", req.ProductNumber, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated succesfully",
	})
}

func (h *Handler) ProductDeleteHandler(c *gin.Context) {
	productId := c.Param("productNumber")
	uid := c.GetInt("user_id")
	companyID := c.GetInt("company_id")

	prod, err := h.service.ReadProduct(productId, companyID)

	if err != nil {
		slog.Error("delete_product - read failed", "error", err, "product_number", productId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = h.service.DeleteProduct(prod.ProductNumber, companyID)
	if err != nil {
		slog.Error("delete_product - failed", "error", err, "product_number", prod.ProductNumber, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("delete_product - success", "product_number", prod.ProductNumber, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted succesfully",
	})
}
