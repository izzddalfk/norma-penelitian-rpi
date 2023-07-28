package entity

import (
	"fmt"

	"gopkg.in/validator.v2"
)

type Goods struct {
	ID     int
	Name   string
	Stocks int
	Price  float64
}

type GoodsConfig struct {
	ID     int    `validate:"nonzero"`
	Name   string `validate:"nonzero"`
	Stocks int
	Price  float64 `validate:"nonzero"`
}

func (g *Goods) IncreaseStock(total int) {
	g.Stocks += total
}

func (g *Goods) DecreaseStock(total int) {
	g.Stocks -= total
}

func NewGoods(cfg GoodsConfig) (*Goods, error) {
	if err := validator.Validate(cfg); err != nil {
		return nil, fmt.Errorf("unable to create goods entity due: %w", err)
	}

	stocks := 0
	if cfg.Stocks > 0 {
		stocks = cfg.Stocks
	}
	goods := &Goods{
		ID:     cfg.ID,
		Name:   cfg.Name,
		Stocks: stocks,
		Price:  cfg.Price,
	}

	return goods, nil
}
