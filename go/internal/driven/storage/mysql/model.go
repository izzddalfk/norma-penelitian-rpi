package storagemysql

import "github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/core/entity"

type GoodsRow struct {
	ID     int     `db:"id"`
	Name   string  `db:"name"`
	Stocks int     `db:"stocks"`
	Price  float64 `db:"price"`
}

func (r GoodsRow) ToGoodsEntity() entity.Goods {
	return entity.Goods(r)
}

type GoodsRowCollection []GoodsRow

func (c GoodsRowCollection) ToGoodsEntityCollection() []entity.Goods {
	var goodsEntityCollection []entity.Goods
	for _, goodsRow := range c {
		goodsEntityCollection = append(goodsEntityCollection, goodsRow.ToGoodsEntity())
	}
	return goodsEntityCollection
}

type ShoppingCartRow struct {
	ID          int64   `db:"id"`
	UserID      int     `db:"id_user"`
	TotalAmount float64 `db:"total_amount"`
	Status      int     `db:"status"`
}

type TransactionRow struct {
	ID          int64   `db:"id"`
	UserID      int     `db:"id_user"`
	TotalAmount float64 `db:"total_amount"`
	Status      int     `db:"status"`
	GoodsID     int     `db:"id_goods"`
	GoodsPrice  float64 `db:"price"`
	TotalGoods  int     `db:"total_goods"`
	CreatedAt   int64   `db:"created_at"`
}

type TransactionRowCollection []TransactionRow

func (r TransactionRowCollection) ToShoppingCartEntity() *entity.ShoppingCart {
	cart := &entity.ShoppingCart{
		ID:          r[0].ID,
		UserID:      r[0].UserID,
		TotalAmount: r[0].TotalAmount,
	}

	for _, trxRow := range r {
		cart.Details = append(cart.Details, entity.ShoppingCartDetail{
			GoodsID:    trxRow.GoodsID,
			GoodsPrice: trxRow.GoodsPrice,
			TotalGoods: trxRow.TotalGoods,
			CreatedAt:  trxRow.CreatedAt,
		})
	}

	return cart
}
