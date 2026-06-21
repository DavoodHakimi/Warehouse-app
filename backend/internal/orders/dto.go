package orders

type OrdersInfo struct {
	Orders []OrderInfoResponse `json:"orders"`
}
type OrderInfoResponse struct {
	ID                  uint           `json:"id"`
	OrderType           string         `json:"order_type"`
	OrderNumber         string         `json:"order_number"`
	Status              string         `json:"status"`
	BusinessPartnerName string         `json:"business_partner_name"`
	Currency            string         `json:"currency"`
	ExchangeRate        float64        `json:"exchange_rate"`
	OrderItems          []OrderItemReq `json:"order_items"`
}
type CreateOrderRequest struct {
	OrderType             string         `binding:"required" json:"order_type"`
	BusinessPartnerID     uint           `binding:"required,numeric" json:"business_partner_name"`
	BusinessPartnerTypeId uint           `binding:"required,numeric" json:"business_partner_type_id"`
	CurrencyID            uint           `binding:"required,numeric" json:"currency"`
	ExchangeRate          float64        `binding:"required" json:"exchange_rate"`
	OrderItems            []OrderItemReq `binding:"required" json:"order_items"`
}

type OrderItemReq struct {
	ProductID    uint    `binding:"required" json:"product_id"`
	Quantity     int     `binding:"required" json:"quantity"`
	PerItemPrice float64 `binding:"required" json:"per_item_price"`
}
type UpdateOrderRequest struct {
	ID                uint    `binding:"required" json:"id"`
	OrderNumber       string  `json:"order_number"`
	OrderType         string  `binding:"required" json:"order_type"`
	BusinessPartnerID uint    `binding:"required,numeric" json:"business_partner_name"`
	CurrencyID        uint    `binding:"required,numeric" json:"currency"`
	ExchangeRate      float64 `binding:"required" json:"exchange_rate"`
}
