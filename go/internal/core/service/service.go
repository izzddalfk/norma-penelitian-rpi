package service

import (
	"context"
	"fmt"

	"github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/core/entity"
	"gopkg.in/validator.v2"
)

type Service interface {
	// small UMKM
	ShowListOfGoods(ctx context.Context, input ShowListOfGoodsInput) ([]entity.Goods, error)
	AddToCart(ctx context.Context, input AddToCartInput) (*AddToCartOutput, error)
	Pay(ctx context.Context, input PayInput) (*entity.Transaction, error)
	// huge UMKM
	ReqCalculateDeliveryPrice(ctx context.Context, input ReqCalculateDeliveryPriceInput) error
	ReqPickupDelivery(ctx context.Context, transactionID int) error
	UpdateStock(ctx context.Context, input UpdateStockInput) error
}

type Storage interface {
	GetGoods(ctx context.Context, input GetGoodsInput) ([]entity.Goods, error)
	GetExistingShoppingCart(ctx context.Context, shoppingCartID int64) (*entity.ShoppingCart, error)
	AddGoodToCart(ctx context.Context, shoppingCart *entity.ShoppingCart) (*entity.ShoppingCart, error)
	CreateTransaction(ctx context.Context, shoppingCart *entity.ShoppingCart) (*entity.Transaction, error)
}

type SupportService interface {
	CalculateDeliveryPrice(ctx context.Context) (float64, error)
	PickupDelivery(ctx context.Context) (bool, error)
}

type service struct {
	storage        Storage
	supportService SupportService
}

type ServiceConfig struct {
	Storage        Storage        `validate:"nonnil"`
	SupportService SupportService `validate:"nonnil"`
}

func NewService(config ServiceConfig) (Service, error) {
	if err := validator.Validate(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &service{
		storage:        config.Storage,
		supportService: config.SupportService,
	}, nil
}

func (s *service) ShowListOfGoods(ctx context.Context, input ShowListOfGoodsInput) ([]entity.Goods, error) {
	goods, err := s.storage.GetGoods(ctx, input.ToGetGoodsStorageInput())
	if err != nil {
		return nil, fmt.Errorf("unable to get list of goods due: %w", err)
	}

	return goods, nil
}

func (s *service) AddToCart(ctx context.Context, input AddToCartInput) (*AddToCartOutput, error) {
	// if there's shopping cart ID in the input, then just get it and update the existing cart
	var shoppingCart entity.ShoppingCart
	switch input.CartID > 0 {
	case true:
		existShoppingCart, err := s.storage.GetExistingShoppingCart(ctx, input.CartID)
		if err != nil {
			return nil, fmt.Errorf("unable to add goods to shopping cart due: %w", err)
		}
		existShoppingCart.AddGoods(entity.AddGoodsInput{
			GoodsID:    input.GoodsID,
			GoodsPrice: input.GoodsPrice,
			TotalGoods: input.Total,
		})
		shoppingCart = *existShoppingCart
	default:
		newShoppingCart, err := entity.NewShoppingCart(entity.ShoppingCartConfig{
			UserID: input.UserID,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to add goods to shopping cart due: %w", err)
		}
		newShoppingCart.AddGoods(entity.AddGoodsInput{
			GoodsID:    input.GoodsID,
			GoodsPrice: input.GoodsPrice,
			TotalGoods: input.Total,
		})
		shoppingCart = *newShoppingCart
	}

	simpleCart, err := s.storage.AddGoodToCart(ctx, &shoppingCart)
	if err != nil {
		return nil, fmt.Errorf("unable to store shopping cart info to storage due: %w", err)
	}

	return &AddToCartOutput{
		CartID:      simpleCart.ID,
		UserID:      simpleCart.UserID,
		TotalGoods:  shoppingCart.GetTotalGoods(),
		TotalAmount: simpleCart.TotalAmount,
	}, nil
}

func (s *service) Pay(ctx context.Context, input PayInput) (*entity.Transaction, error) {
	existShoppingCart, err := s.storage.GetExistingShoppingCart(ctx, input.CartID)
	if err != nil {
		return nil, fmt.Errorf("unable to get shopping cart information due: %w", err)
	}

	return s.storage.CreateTransaction(ctx, existShoppingCart)
}

func (s *service) ReqCalculateDeliveryPrice(ctx context.Context, input ReqCalculateDeliveryPriceInput) error {
	return nil
}

func (s *service) ReqPickupDelivery(ctx context.Context, transactionID int) error {
	return nil
}

func (s *service) UpdateStock(ctx context.Context, input UpdateStockInput) error {
	return nil
}
