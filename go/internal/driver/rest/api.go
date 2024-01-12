package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/core/service"
	"gopkg.in/validator.v2"
)

type api struct {
	id     string
	servce service.Service
}

type APIConfig struct {
	Service service.Service `validate:"nonnil"`
}

func NewAPI(config APIConfig) (*api, error) {
	if err := validator.Validate(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	svcID := uuid.NewString()

	return &api{
		id:     svcID,
		servce: config.Service,
	}, nil
}

func (a *api) Handler() *gin.Engine {
	r := gin.Default()
	// small umkm API
	smallRouter := r.Group("/api/small")
	{
		smallRouter.GET("/stocks", a.HandleShowListOfGoods)
		smallRouter.POST("/cart", a.HandleAddGoodsToCart)
		smallRouter.POST("/pay", a.HandlePay)
	}
	// for testing API
	r.POST("/clear-db", a.HandleClearDB)

	return r
}

func (a *api) HandleShowListOfGoods(c *gin.Context) {
	var qpErrors []string
	qpPage, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		qpErrors = append(qpErrors, err.Error())
	}
	qpTotalGoods, err := strconv.Atoi(c.DefaultQuery("total_goods", "3"))
	if err != nil {
		qpErrors = append(qpErrors, err.Error())
	}
	if len(qpErrors) > 0 {
		c.JSON(
			http.StatusBadRequest,
			NewBadRequestErrorResponse(qpErrors),
		)
		return
	}

	// call main service
	listOfGoods, err := a.servce.ShowListOfGoods(c.Request.Context(), service.ShowListOfGoodsInput{
		Page:       qpPage,
		TotalGoods: qpTotalGoods,
		Sort:       service.Sort(c.DefaultQuery("sort", "DESC")),
		SortBy:     c.DefaultQuery("sort_by", "id"),
	})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			NewInternalServerErrorResponse(err.Error()),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		NewSuccessResponse(listOfGoods, a.id),
	)
}

func (a *api) HandleAddGoodsToCart(c *gin.Context) {
	var reqBody struct {
		CartID     int     `json:"cart_id"`
		UserID     int     `json:"user_id" binding:"required"`
		GoodsID    int     `json:"goods_id" binding:"required"`
		GoodsPrice float64 `json:"goods_price" binding:"required"`
		TotalGoods int     `json:"total_goods" binding:"required"`
	}

	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			NewBadRequestErrorResponse(err.Error()),
		)
		return
	}

	output, err := a.servce.AddToCart(c.Request.Context(), service.AddToCartInput{
		CartID:     int64(reqBody.CartID),
		UserID:     reqBody.UserID,
		GoodsID:    reqBody.GoodsID,
		GoodsPrice: reqBody.GoodsPrice,
		Total:      reqBody.TotalGoods,
	})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			NewInternalServerErrorResponse(err.Error()),
		)
		return
	}

	var respBody struct {
		CartID      int64   `json:"cart_id"`
		TotalGoods  int     `json:"total_goods"`
		TotalAmount float64 `json:"total_amount"`
	}
	respBody.CartID = output.CartID
	respBody.TotalGoods = output.TotalGoods
	respBody.TotalAmount = output.TotalAmount

	c.JSON(http.StatusOK, NewSuccessResponse(respBody, a.id))
}

func (a *api) HandlePay(c *gin.Context) {
	var reqBody struct {
		CartID        int64   `json:"cart_id" binding:"required"`
		PaymentAmount float64 `json:"payment_amount" binding:"required"`
	}

	err := c.ShouldBindJSON(&reqBody)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			NewBadRequestErrorResponse(err.Error()),
		)
		return
	}

	trx, err := a.servce.Pay(c.Request.Context(), service.PayInput{
		CartID:        reqBody.CartID,
		PaymentAmount: reqBody.PaymentAmount,
	})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			NewInternalServerErrorResponse(err.Error()),
		)
		return
	}

	var respBody struct {
		TransactionID int64   `json:"transaction_id"`
		TotalAmount   float64 `json:"total_amount"`
		PaymentAmount float64 `json:"payment_amount"`
		ReturnAmount  float64 `json:"return_amount"`
	}
	respBody.TransactionID = trx.ID
	respBody.TotalAmount = trx.TotalAmount
	respBody.PaymentAmount = trx.PaymentAmount
	respBody.ReturnAmount = trx.ReturnAmount

	c.JSON(http.StatusOK, NewSuccessResponse(respBody, a.id))
}

func (a *api) HandleClearDB(c *gin.Context) {
	if err := a.servce.ClearDatabase(c.Request.Context()); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			NewInternalServerErrorResponse(err.Error()),
		)
		return
	}

	c.JSON(http.StatusOK, NewSuccessResponse("Clear database success!", a.id))
}
