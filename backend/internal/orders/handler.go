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
	h.service.AllOrders(userCompanyID.(int))
}

func (h *Handler) OrderInfoHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	h.service.ReadOrder(orderId)

}

func (h *Handler) OrderCreationHandler(c *gin.Context) {

	h.service.CreateOrder()
}

func (h *Handler) OrderUpdateHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	h.service.UpdateOrder(&UpdateOrderRequest{}, orderId)

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
	h.service.DeleteOrder(order.ID)

}

func (h *Handler) OrderApproveHandler(c *gin.Context) {
	orderId := c.Param("orderID")

}
func (h *Handler) OrderPackHandler(c *gin.Context) {
	orderId := c.Param("orderID")

}
func (h *Handler) OrderShipHandler(c *gin.Context) {
	orderId := c.Param("orderID")

}
func (h *Handler) OrderReceiveHandler(c *gin.Context) {
	orderId := c.Param("orderID")

}
