package service

type UpdateStockAction string

type Sort string

const (
	IncreaseStock UpdateStockAction = "INCR"
	DecreaseStock UpdateStockAction = "DECR"
	SortAsc       Sort              = "ASC"
	SortDesc      Sort              = "DESC"
)

type ShowListOfGoodsInput struct {
	Page       int
	TotalGoods int
	Sort       Sort
	SortBy     string
}

func (i ShowListOfGoodsInput) ToGetGoodsStorageInput() GetGoodsInput {
	// default values
	input := GetGoodsInput{
		Offset: 0,
		Limit:  10,
		Sort:   SortDesc,
		SortBy: "id",
	}
	if i.Page > 0 {
		input.Offset = i.Page - 1
	}
	if i.TotalGoods > 0 {
		input.Limit = i.TotalGoods
	}
	if len(i.Sort) > 0 {
		input.Sort = i.Sort
	}
	if len(i.SortBy) > 0 {
		input.SortBy = i.SortBy
	}
	return input
}

type AddToCartInput struct {
	CartID     int64
	UserID     int
	GoodsID    int
	GoodsPrice float64
	Total      int
}

type AddToCartOutput struct {
	CartID      int64
	UserID      int
	TotalGoods  int
	TotalAmount float64
}

type PayInput struct {
	CartID int64
}

type ReqCalculateDeliveryPriceInput struct {
	Location  float64
	GoodsSpec GoodsSpecification
}

type UpdateStockInput struct {
	Action  UpdateStockAction
	GoodsID int
	Total   int
}

type GoodsSpecification struct {
	Weight float32
	Volume float32
}

type GetGoodsInput struct {
	Offset int
	Limit  int
	Sort   Sort
	SortBy string
}

type CreateTransactionInput struct {
	CartID      int
	TotalAmount float64
}
