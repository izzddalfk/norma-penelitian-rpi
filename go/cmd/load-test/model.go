package main

type respBodyAddToCart struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		CartID      int     `json:"cart_id"`
		TotalAmount float64 `json:"total_amount"`
	} `json:"data,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
}

type addToCartReqBody struct {
	CartID     int     `json:"cart_id,omitempty"`
	UserID     int     `json:"user_id"`
	GoodsID    int     `json:"goods_id"`
	GoodsPrice float64 `json:"goods_price"`
	TotalGoods int     `json:"total_goods"`
}

type payReqBody struct {
	CartID      int     `json:"cart_id"`
	TotalAmount float64 `json:"total_amount"`
}
