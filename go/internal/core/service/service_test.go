package service_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/core/entity"
	"github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/core/service"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
)

func TestNewService(mainT *testing.T) {
	deps := newMockDependencies(mockDependenciesConfig{})

	testCases := []struct {
		Name    string
		Config  service.ServiceConfig
		IsError bool
	}{
		{
			Name: "Test valid config",
			Config: service.ServiceConfig{
				Storage:        deps.Storage,
				SupportService: deps.SupportService,
			},
			IsError: false,
		},
		{
			Name: "Test missing storage",
			Config: service.ServiceConfig{
				Storage:        nil,
				SupportService: deps.SupportService,
			},
			IsError: true,
		},
		{
			Name: "Test missing support service",
			Config: service.ServiceConfig{
				Storage:        deps.Storage,
				SupportService: nil,
			},
			IsError: true,
		},
	}

	for _, testCase := range testCases {
		mainT.Run(testCase.Name, func(t *testing.T) {
			_, err := service.NewService(testCase.Config)
			require.Equal(t, testCase.IsError, (err != nil), "unexpected error")
		})
	}
}

func TestShowStocks(mainT *testing.T) {
	dummyGoods := newDummyGoods(10)
	testCases := []struct {
		Name                 string
		Input                service.ShowListOfGoodsInput
		DummyGoodsCollection []entity.Goods
		ExpectedGoods        []entity.Goods
	}{
		{
			Name:                 "Test show list of goods with default input",
			Input:                service.ShowListOfGoodsInput{},
			DummyGoodsCollection: dummyGoods,
			ExpectedGoods:        dummyGoods,
		},
		{
			Name: "Test show list of goods with modified input",
			Input: service.ShowListOfGoodsInput{
				Page:       2,
				TotalGoods: 5,
			},
			DummyGoodsCollection: dummyGoods,
			ExpectedGoods:        dummyGoods[5:],
		},
		// TODO: Sort list of goods
	}

	for _, testCase := range testCases {
		mainT.Run(testCase.Name, func(t *testing.T) {
			deps := newMockDependencies(mockDependenciesConfig{
				mockStorageDummyGoods: testCase.DummyGoodsCollection,
			})

			svc, err := service.NewService(service.ServiceConfig(deps))
			require.NoError(t, err)

			goods, err := svc.ShowListOfGoods(context.Background(), testCase.Input)
			require.NoError(t, err)
			require.ElementsMatch(t, testCase.ExpectedGoods, goods)
		})
	}
}

func TestAddToCart(mainT *testing.T) {
	testCases := []struct {
		Name           string
		Input          []service.AddToCartInput
		ExpectedOutput service.AddToCartOutput
	}{
		{
			Name: "Successfully add goods to cart from empty cart",
			Input: []service.AddToCartInput{
				{
					UserID:     100,
					GoodsID:    1,
					GoodsPrice: 2000,
					Total:      2,
				},
			},
			ExpectedOutput: service.AddToCartOutput{
				TotalGoods:  2,
				TotalAmount: 2 * 2000,
			},
		},
		{
			Name: "Successfully add goods to cart from non-empty cart",
			Input: []service.AddToCartInput{
				{
					UserID:     100,
					GoodsID:    2,
					GoodsPrice: 2000,
					Total:      1,
				},
				{
					CartID:     1,
					UserID:     100,
					GoodsID:    3,
					GoodsPrice: 1500,
					Total:      4,
				},
			},
			ExpectedOutput: service.AddToCartOutput{
				TotalGoods:  5,
				TotalAmount: (1 * 2000) + (4 * 1500),
			},
		},
	}

	for _, testCase := range testCases {
		mainT.Run(testCase.Name, func(t *testing.T) {
			deps := newMockDependencies(mockDependenciesConfig{})

			svc, err := service.NewService(service.ServiceConfig(deps))
			require.NoError(t, err)

			actualOutput := service.AddToCartOutput{
				TotalGoods:  0,
				TotalAmount: 0,
			}
			for _, input := range testCase.Input {
				output, err := svc.AddToCart(context.Background(), input)
				require.NoError(t, err)
				actualOutput.TotalGoods = output.TotalGoods
				actualOutput.TotalAmount = output.TotalAmount
			}
			require.Equal(t, testCase.ExpectedOutput, actualOutput)
		})
	}
}

func TestPay(mainT *testing.T) {
	testCases := []struct {
		Name              string
		ShoppingCartInput service.AddToCartInput
		Input             service.PayInput
		ExpectedSuccess   bool
	}{
		{
			Name: "Successfully do payment",
			ShoppingCartInput: service.AddToCartInput{
				UserID:     200,
				GoodsID:    1,
				GoodsPrice: 2000,
				Total:      2,
			},
			Input: service.PayInput{
				CartID: 1,
			},
			ExpectedSuccess: true,
		},
	}

	for _, testCase := range testCases {
		mainT.Run(testCase.Name, func(t *testing.T) {
			deps := newMockDependencies(mockDependenciesConfig{})

			svc, err := service.NewService(service.ServiceConfig(deps))
			require.NoError(t, err)

			ctx := context.Background()

			// add to cart first
			_, err = svc.AddToCart(ctx, testCase.ShoppingCartInput)
			require.NoError(t, err)

			trx, err := svc.Pay(ctx, testCase.Input)
			require.NoError(t, err)
			require.Equal(t, trx.ID != "", testCase.ExpectedSuccess)
		})
	}
}

type mockDependencies struct {
	Storage        service.Storage
	SupportService service.SupportService
}

type mockDependenciesConfig struct {
	mockStorageDummyGoods []entity.Goods
}

func newMockDependencies(config mockDependenciesConfig) mockDependencies {
	return mockDependencies{
		Storage: &mockStorage{
			Goods:        config.mockStorageDummyGoods,
			ShoppingCart: map[int64]entity.ShoppingCart{},
		},
		SupportService: &mockSupportService{},
	}
}

func newDummyGoods(total int) []entity.Goods {
	f := faker.New()
	rg := rand.New(rand.NewSource(time.Now().UnixNano()))
	output := []entity.Goods{}
	for i := 1; i <= total; i++ {
		output = append(output, entity.Goods{
			ID:     i,
			Name:   f.Food().Fruit(),
			Stocks: rg.Int(),
			Price:  rg.Float64(),
		})
	}

	return output
}

type mockStorage struct {
	Goods        []entity.Goods
	ShoppingCart map[int64]entity.ShoppingCart
}

func (m *mockStorage) GetGoods(ctx context.Context, input service.GetGoodsInput) ([]entity.Goods, error) {
	if len(m.Goods) == 0 || m.Goods == nil {
		return []entity.Goods{}, nil
	}

	if input.Limit > len(m.Goods) {
		return m.Goods, nil
	}

	return m.Goods[input.Offset*input.Limit : (input.Offset+1)*input.Limit], nil
}

func (m *mockStorage) GetExistingShoppingCart(ctx context.Context, shoppingCartID int64) (*entity.ShoppingCart, error) {
	existCart, ok := m.ShoppingCart[shoppingCartID]
	if !ok {
		return nil, fmt.Errorf("shopping cart with %d ID not exist yet", shoppingCartID)
	}
	return &existCart, nil
}

func (m *mockStorage) AddGoodToCart(ctx context.Context, cart *entity.ShoppingCart) error {
	if cart.ID > 0 {
		existCart, ok := m.ShoppingCart[cart.ID]
		if !ok {
			return fmt.Errorf("unexpected error: no existing cart for ID %d", cart.ID)
		}
		existCart.Details = append(existCart.Details, cart.Details...)
		existCart.TotalAmount = existCart.GetTotalAmount()

		m.ShoppingCart[cart.ID] = existCart
	} else {
		newCartID := len(m.ShoppingCart) + 1
		m.ShoppingCart[int64(newCartID)] = entity.ShoppingCart{
			ID:          int64(newCartID),
			UserID:      cart.UserID,
			TotalAmount: cart.GetTotalAmount(),
			Details:     cart.Details,
		}
	}
	return nil
}

func (m *mockStorage) CreateTransaction(ctx context.Context, shoppingCart *entity.ShoppingCart) (*entity.Transaction, error) {
	if shoppingCart == nil {
		return nil, fmt.Errorf("no shopping cart")
	}
	return &entity.Transaction{
		ID: faker.New().UUID().V4(),
	}, nil
}

type mockSupportService struct{}

func (m *mockSupportService) CalculateDeliveryPrice(ctx context.Context) (float64, error) {
	return 0, nil
}

func (m *mockSupportService) PickupDelivery(ctx context.Context) (bool, error) {
	return true, nil
}
