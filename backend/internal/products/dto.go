package products

type ProductsInfo struct {
	Products []ProductInfoResponse `json:"products"`
}

type ProductInfoResponse struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	ProductNumber string  `json:"product_number"`
	IsFrozen      bool    `json:"is_frozen"`
	DefaultPrice  float64 `json:"default_price"`
}

type ProductRequest struct {
	Name         string  `json:"name" binding:"required,min=4,max=100"`
	IsFrozen     bool    `json:"is_frozen" binding:"boolean"`
	DefaultPrice float64 `json:"default_price" binding:"numeric"`
}

type UpdateProductRequest struct {
	ID            int     `json:"id" binding:"required,numeric"`
	Name          string  `json:"name" binding:"required,min=4,max=100"`
	ProductNumber string  `json:"product_number" binding:"required"`
	IsFrozen      bool    `json:"is_frozen" binding:"boolean"`
	DefaultPrice  float64 `json:"default_price" binding:"numeric"`
}
