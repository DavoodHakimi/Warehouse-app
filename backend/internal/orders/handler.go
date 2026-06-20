package orders

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

func (h *Handler) AllOrdersHandler(c *gin.Context) {
	userCompanyID, exist := c.Get("company_id")
	uid := c.GetInt("user_id")
	if !exist {
		slog.Error("all_orders - company_id missing", "user_id", uid)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you don't have access to this page",
		})
		return
	}
	allOrders, err := h.service.AllOrders(userCompanyID.(int))
	if err != nil {
		slog.Error("all_orders - failed", "error", err, "company_id", userCompanyID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.Info("all_orders - success", "company_id", userCompanyID, "count", len(allOrders.Orders))
	c.JSON(http.StatusOK, *allOrders)
}

func (h *Handler) OrderInfoHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	uid := c.GetInt("user_id")
	cid := c.GetInt("company_id")

	order, err := h.service.ReadOrder(orderId, cid)

	if err != nil {
		slog.Error("order_info - failed", "error", err, "order_id", orderId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("order_info - success", "order_id", orderId, "request_by", uid)
	c.JSON(http.StatusOK, *order)

}

func (h *Handler) OrderCreationHandler(c *gin.Context) {
	var req CreateOrderRequest
	uid := c.GetInt("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("create_order - validation failed", "error", err, "request_by", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}
	comapnyID := c.GetInt("company_id")
	err := h.service.CreateOrder(&req, comapnyID)
	if err != nil {
		slog.Error("create_order - failed", "error", err, "order_type", req.OrderType, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("create_order - success", "order_type", req.OrderType, "partner_id", req.BusinessPartnerID, "request_by", uid)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created succesfully",
	})
}

func (h *Handler) OrderUpdateHandler(c *gin.Context) {
	var req UpdateOrderRequest
	uid := c.GetInt("user_id")
	cid := c.GetInt("company_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("update_order - validation failed", "error", err, "request_by", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation Failed.",
			"details": err.Error(),
		})
		return
	}
	err := h.service.UpdateOrder(&req, uid, cid)
	if err != nil {
		slog.Error("update_order - failed", "error", err, "order_id", req.ID, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	slog.Info("update_order - success", "order_id", req.ID, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order updated succesfully",
	})

}

func (h *Handler) OrderDeleteHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	uid := c.GetInt("user_id")
	cid := c.GetInt("company_id")

	order, err := h.service.ReadOrder(orderId, cid)
	if err != nil {
		slog.Error("delete_order - read failed", "error", err, "order_id", orderId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = h.service.DeleteOrder(order.ID, cid)

	if err != nil {
		slog.Error("delete_order - failed", "error", err, "order_id", order.ID, "request_by", uid)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("delete_order - success", "order_id", order.ID, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order deleted succesfully",
	})

}

func (h *Handler) OrderApproveHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	uid := c.GetInt("user_id")
	cid := c.GetInt("company_id")

	err := h.service.Approve(orderId, cid)
	if err != nil {
		slog.Error("order_approve - failed", "error", err, "order_id", orderId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("order_approve - success", "order_id", orderId, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Approved Successfully",
	})

}
func (h *Handler) OrderPackHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	uid := c.GetInt("user_id")
	cid := c.GetInt("company_id")

	err := h.service.Pack(orderId, cid)
	if err != nil {
		slog.Error("order_pack - failed", "error", err, "order_id", orderId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("order_pack - success", "order_id", orderId, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Packed Successfully",
	})
}
func (h *Handler) OrderShipHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	uid := c.GetInt("user_id")
	cid := c.GetInt("company_id")

	err := h.service.Ship(orderId, cid)
	if err != nil {
		slog.Error("order_ship - failed", "error", err, "order_id", orderId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("order_ship - success", "order_id", orderId, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Shipped Successfully",
	})

}
func (h *Handler) OrderReceiveHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	uid := c.GetInt("user_id")
	cid := c.GetInt("company_id")

	err := h.service.Receive(orderId, cid)
	if err != nil {
		slog.Error("order_receive - failed", "error", err, "order_id", orderId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("order_receive - success", "order_id", orderId, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Received Successfully",
	})
}

func (h *Handler) OrderWaitingHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	uid := c.GetInt("user_id")
	cid := c.GetInt("company_id")

	err := h.service.MarkWaiting(orderId, cid)
	if err != nil {
		slog.Error("order_wait - failed", "error", err, "order_id", orderId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("order_wait - success", "order_id", orderId, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Waiting Successfully",
	})
}

func (h *Handler) OrderCancelHandler(c *gin.Context) {
	orderId := c.Param("orderID")
	uid := c.GetInt("user_id")
	cid := c.GetInt("company_id")

	err := h.service.Cancel(orderId, cid)
	if err != nil {
		slog.Error("order_cancel - failed", "error", err, "order_id", orderId, "request_by", uid)
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	slog.Info("order_cancel - success", "order_id", orderId, "request_by", uid)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order status changed to Canceled Successfully",
	})
}
