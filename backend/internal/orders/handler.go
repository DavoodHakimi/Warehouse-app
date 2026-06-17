package orders

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

func (h *Handler) AllOrdersHandler(c *gin.Context) {
	userCompanyID, exist := c.Get("company_id")
	if !exist {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you don't have access to this page",
		})
		return
	}
	allOrders, err := h.service.AllOrders(userCompanyID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *allOrders)
}

func (h *Handler) OrderInfoHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	order, err := h.service.ReadOrder(orderId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *order)

}

func (h *Handler) OrderCreationHandler(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
	}
	comapnyID, _ := c.Get("company_id")
	err := h.service.CreateOrder(&req, comapnyID.(int))
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

func (h *Handler) OrderUpdateHandler(c *gin.Context) {
	var req UpdateOrderRequest
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
	err := h.service.UpdateOrder(&req, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Order updated succesfully",
	})

}

func (h *Handler) OrderDeleteHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	order, err := h.service.ReadOrder(orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = h.service.DeleteOrder(order.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order deleted succesfully",
	})

}

func (h *Handler) OrderApproveHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	err := h.service.Approve(orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Approved Successfully",
	})

}
func (h *Handler) OrderPackHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	err := h.service.Pack(orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Packed Successfully",
	})
}
func (h *Handler) OrderShipHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	err := h.service.Ship(orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Shipped Successfully",
	})

}
func (h *Handler) OrderReceiveHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	err := h.service.Receive(orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Received Successfully",
	})
}

func (h *Handler) OrderWaitingHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	err := h.service.MarkWaiting(orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Waiting Successfully",
	})
}

func (h *Handler) OrderCancelHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	err := h.service.Cancel(orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Canceled Successfully",
	})
}
