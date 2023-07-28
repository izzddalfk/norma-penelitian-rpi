package storagemysql_test

import (
	"context"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/core/entity"
	"github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/core/service"
	storagemysql "github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/driven/storage/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestGetGoods(mainT *testing.T) {
	dbConn := initDB(mainT)
	defer dbConn.Close()

	strg, err := storagemysql.NewStorage(storagemysql.StorageConfig{
		DBClient: dbConn,
	})
	require.NoError(mainT, err)

	listOfGoods, err := strg.GetGoods(context.Background(), service.GetGoodsInput{
		Offset: 0,
		Limit:  5,
		Sort:   service.SortDesc,
		SortBy: "id",
	})
	require.NoError(mainT, err)
	require.NotZero(mainT, listOfGoods)
}

func TestGetExistingShoppingCart(mainT *testing.T) {
	dbConn := initDB(mainT)
	defer func() {
		// clean transactions and transaction details table
		dbConn.ExecContext(context.Background(), "TRUNCATE transactions")
		dbConn.ExecContext(context.Background(), "TRUNCATE transaction_details")
		dbConn.Close()
	}()

	strg, err := storagemysql.NewStorage(storagemysql.StorageConfig{
		DBClient: dbConn,
	})
	require.NoError(mainT, err)

	nilCart, err := strg.GetExistingShoppingCart(context.Background(), 1)
	require.NoError(mainT, err)
	require.Nil(mainT, nilCart)

	// create new shopping cart
	newShoppingCartInput := &entity.ShoppingCart{
		UserID: 100,
		Details: []entity.ShoppingCartDetail{
			{
				GoodsID:    1,
				TotalGoods: 1,
				GoodsPrice: 3000,
				CreatedAt:  1689873350,
			},
			{
				GoodsID:    2,
				TotalGoods: 3,
				GoodsPrice: 1500,
				CreatedAt:  1689873350,
			},
		},
	}
	_, err = strg.AddGoodToCart(context.Background(), newShoppingCartInput)
	require.NoError(mainT, err)

	// assert the existing cart
	expectedCart := &entity.ShoppingCart{
		ID:          1,
		UserID:      100,
		TotalAmount: 3000 + (3 * 1500),
		Details: []entity.ShoppingCartDetail{
			{
				GoodsID:    1,
				TotalGoods: 1,
				GoodsPrice: 3000,
				CreatedAt:  1689873350,
			},
			{
				GoodsID:    2,
				TotalGoods: 3,
				GoodsPrice: 1500,
				CreatedAt:  1689873350,
			},
		},
	}
	existingCart, err := strg.GetExistingShoppingCart(context.Background(), 1)
	require.NoError(mainT, err)
	require.Equal(mainT, expectedCart, existingCart)
}

func TestAddGoodToCart(mainT *testing.T) {
	dbConn := initDB(mainT)
	defer dbConn.Close()

	strg, err := storagemysql.NewStorage(storagemysql.StorageConfig{
		DBClient: dbConn,
	})
	require.NoError(mainT, err)

	testCases := []struct {
		Name          string
		Input         []entity.ShoppingCart
		ExpectedCarts []entity.ShoppingCart
	}{
		{
			Name: "Shopping cart empty, add new shopping cart",
			Input: []entity.ShoppingCart{
				{
					UserID: 100,
					Details: []entity.ShoppingCartDetail{
						{
							GoodsID:    1,
							TotalGoods: 1,
							GoodsPrice: 3000,
							CreatedAt:  1689873350,
						},
						{
							GoodsID:    2,
							TotalGoods: 3,
							GoodsPrice: 1500,
							CreatedAt:  1689873350,
						},
					},
				},
			},
			ExpectedCarts: []entity.ShoppingCart{
				{
					ID:          1,
					UserID:      100,
					TotalAmount: 3000 + (3 * 1500),
					Details: []entity.ShoppingCartDetail{
						{
							GoodsID:    1,
							TotalGoods: 1,
							GoodsPrice: 3000,
							CreatedAt:  1689873350,
						},
						{
							GoodsID:    2,
							TotalGoods: 3,
							GoodsPrice: 1500,
							CreatedAt:  1689873350,
						},
					},
				},
			},
		},
		{
			Name: "Add goods to existing shopping cart",
			Input: []entity.ShoppingCart{
				{
					UserID: 100,
					Details: []entity.ShoppingCartDetail{
						{
							GoodsID:    1,
							TotalGoods: 1,
							GoodsPrice: 3000,
							CreatedAt:  1689873350,
						},
						{
							GoodsID:    2,
							TotalGoods: 3,
							GoodsPrice: 1500,
							CreatedAt:  1689873350,
						},
					},
				},
				{
					ID:     1,
					UserID: 100,
					Details: []entity.ShoppingCartDetail{
						{
							GoodsID:    6,
							TotalGoods: 1,
							GoodsPrice: 2500,
							CreatedAt:  1689873500,
						},
					},
				},
			},
			ExpectedCarts: []entity.ShoppingCart{
				{
					ID:          1,
					UserID:      100,
					TotalAmount: 3000 + (3 * 1500) + 2500,
					Details: []entity.ShoppingCartDetail{
						{
							GoodsID:    1,
							TotalGoods: 1,
							GoodsPrice: 3000,
							CreatedAt:  1689873350,
						},
						{
							GoodsID:    2,
							TotalGoods: 3,
							GoodsPrice: 1500,
							CreatedAt:  1689873350,
						},
						{
							GoodsID:    6,
							TotalGoods: 1,
							GoodsPrice: 2500,
							CreatedAt:  1689873500,
						},
					},
				},
			},
		},
		{
			Name: "Shopping cart empty, add new two different shopping carts",
			Input: []entity.ShoppingCart{
				{
					UserID: 100,
					Details: []entity.ShoppingCartDetail{
						{
							GoodsID:    1,
							TotalGoods: 1,
							GoodsPrice: 3000,
							CreatedAt:  1689873350,
						},
						{
							GoodsID:    2,
							TotalGoods: 3,
							GoodsPrice: 1500,
							CreatedAt:  1689873350,
						},
					},
				},
				{
					UserID: 200,
					Details: []entity.ShoppingCartDetail{
						{
							GoodsID:    6,
							TotalGoods: 1,
							GoodsPrice: 2500,
							CreatedAt:  1689873500,
						},
					},
				},
			},
			ExpectedCarts: []entity.ShoppingCart{
				{
					ID:          1,
					UserID:      100,
					TotalAmount: 3000 + (3 * 1500),
					Details: []entity.ShoppingCartDetail{
						{
							GoodsID:    1,
							TotalGoods: 1,
							GoodsPrice: 3000,
							CreatedAt:  1689873350,
						},
						{
							GoodsID:    2,
							TotalGoods: 3,
							GoodsPrice: 1500,
							CreatedAt:  1689873350,
						},
					},
				},
				{
					ID:          2,
					UserID:      200,
					TotalAmount: 2500,
					Details: []entity.ShoppingCartDetail{
						{
							GoodsID:    6,
							TotalGoods: 1,
							GoodsPrice: 2500,
							CreatedAt:  1689873500,
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		mainT.Run(testCase.Name, func(t *testing.T) {
			defer func() {
				// clean transactions and transaction details table
				dbConn.ExecContext(context.Background(), "TRUNCATE transactions")
				dbConn.ExecContext(context.Background(), "TRUNCATE transaction_details")
			}()

			for _, input := range testCase.Input {
				cartOutput, err := strg.AddGoodToCart(context.Background(), &input)
				require.NoError(t, err)
				require.NotNil(t, cartOutput)
			}

			for _, expectedCart := range testCase.ExpectedCarts {
				existingCart, err := strg.GetExistingShoppingCart(context.Background(), expectedCart.ID)
				require.NoError(mainT, err)
				require.Equal(mainT, &expectedCart, existingCart)
			}
		})
	}
}

func TestCreateTransaction(mainT *testing.T) {
	dbConn := initDB(mainT)
	defer func() {
		// clean transactions and transaction details table
		dbConn.ExecContext(context.Background(), "TRUNCATE transactions")
		dbConn.ExecContext(context.Background(), "TRUNCATE transaction_details")
		dbConn.Close()
	}()

	strg, err := storagemysql.NewStorage(storagemysql.StorageConfig{
		DBClient: dbConn,
	})
	require.NoError(mainT, err)

	// create new shopping cart first
	newShoppingCartInput := &entity.ShoppingCart{
		UserID: 100,
		Details: []entity.ShoppingCartDetail{
			{
				GoodsID:    1,
				TotalGoods: 1,
				GoodsPrice: 3000,
				CreatedAt:  1689873350,
			},
			{
				GoodsID:    2,
				TotalGoods: 3,
				GoodsPrice: 1500,
				CreatedAt:  1689873350,
			},
		},
	}
	cartOutput, err := strg.AddGoodToCart(context.Background(), newShoppingCartInput)
	require.NoError(mainT, err)
	require.NotNil(mainT, cartOutput)

	// create new transaction
	newTrx, err := strg.CreateTransaction(context.Background(), &entity.ShoppingCart{ID: 1})
	require.NoError(mainT, err)
	require.Equal(mainT, int64(1), newTrx.ID)
}

func initDB(mainT *testing.T) *sqlx.DB {
	ctx := context.Background()
	sqlDSN := os.Getenv("DB_SQLDSN")
	if len(sqlDSN) == 0 {
		sqlDSN = "root:test1234@tcp(localhost:23306)/umkm?timeout=5s"
	}

	dbConn, err := sqlx.ConnectContext(ctx, "mysql", sqlDSN)
	require.NoError(mainT, err)
	require.NotNil(mainT, dbConn)

	return dbConn
}
