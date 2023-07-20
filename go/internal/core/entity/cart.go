package entity

import (
	"fmt"
	"time"

	"gopkg.in/validator.v2"
)

type ShoppingCart struct {
	ID          int64
	UserID      int
	TotalAmount float64
	Details     []ShoppingCartDetail
}

type ShoppingCartConfig struct {
	UserID int `validate:"nonzero"`
}

func NewShoppingCart(config ShoppingCartConfig) (*ShoppingCart, error) {
	if err := validator.Validate(config); err != nil {
		return nil, fmt.Errorf("unable to create new shopping cart due: %w", err)
	}

	return &ShoppingCart{
		UserID: config.UserID,
	}, nil
}

type AddGoodsInput struct {
	GoodsID    int     `validate:"nonzero"`
	TotalGoods int     `validate:"nonzero"`
	GoodsPrice float64 `validate:"nonzero"`
}

func (c *ShoppingCart) AddGoods(input AddGoodsInput) error {
	if err := validator.Validate(input); err != nil {
		return fmt.Errorf("unable to add goods into cart due: %w", err)
	}

	c.Details = append(c.Details, ShoppingCartDetail{
		GoodsID:    input.GoodsID,
		TotalGoods: input.TotalGoods,
		GoodsPrice: input.GoodsPrice,
		CreatedAt:  time.Now().Unix(),
	})
	// update the shopping cart total amount as well
	c.TotalAmount = c.GetTotalAmount()

	return nil
}

func (c ShoppingCart) GetTotalGoods() int {
	var totalGoods int
	for _, detail := range c.Details {
		totalGoods += detail.TotalGoods
	}
	return totalGoods
}

func (c *ShoppingCart) GetTotalAmount() float64 {
	var newAmount float64
	for _, detail := range c.Details {
		newAmount += float64(detail.TotalGoods) * detail.GoodsPrice
	}

	return newAmount
}

type ShoppingCartDetail struct {
	GoodsID    int
	TotalGoods int
	GoodsPrice float64
	CreatedAt  int64
}
