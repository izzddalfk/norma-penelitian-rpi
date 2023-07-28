package storagemysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/core/entity"
	"github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/core/service"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

type storage struct {
	client *sqlx.DB
}

type StorageConfig struct {
	DBClient *sqlx.DB `validate:"nonnil"`
}

func NewStorage(config StorageConfig) (*storage, error) {
	if err := validator.Validate(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &storage{
		client: config.DBClient,
	}, nil
}

func (s *storage) GetGoods(ctx context.Context, input service.GetGoodsInput) ([]entity.Goods, error) {
	query := fmt.Sprintf(`
			SELECT 
				* 
			FROM goods
			ORDER BY '%s' %s
			LIMIT %d OFFSET %d
		`,
		input.SortBy,
		input.Sort,
		input.Limit,
		input.Offset,
	)

	var goodsCollection GoodsRowCollection
	err := s.client.SelectContext(
		ctx,
		&goodsCollection,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to execute select query for goods due: %w", err)
	}

	return goodsCollection.ToGoodsEntityCollection(), nil
}

func (s *storage) GetExistingShoppingCart(ctx context.Context, shoppingCartID int64) (*entity.ShoppingCart, error) {
	query := `
		SELECT 
			trx.*,
			trx_details.id_goods,
			trx_details.total_goods,
			trx_details.created_at,
			g.price
		FROM transactions trx 
		JOIN transaction_details trx_details
			ON trx.id = trx_details.id_transaction
		JOIN goods g
			ON trx_details.id_goods = g.id
		WHERE trx.id = ? AND trx.status = 0
	`

	var existingCart TransactionRowCollection
	err := s.client.SelectContext(
		ctx,
		&existingCart,
		query,
		shoppingCartID,
	)
	if len(existingCart) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("unable to execute select query for get existing cart due: %w", err)
	}

	return existingCart.ToShoppingCartEntity(), nil
}

func (s *storage) AddGoodToCart(ctx context.Context, shoppingCart *entity.ShoppingCart) (*entity.ShoppingCart, error) {
	dbTx, err := s.client.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction for add goods to cart query: %w", err)
	}
	// Defer the transaction’s rollback. If the transaction succeeds, it will be committed before the function exits,
	// making the deferred rollback call a no-op.
	// If the transaction fails it won’t be committed, meaning that the rollback will be called as the function exits.
	defer dbTx.Rollback()

	simpleCart := &entity.ShoppingCart{
		UserID: shoppingCart.UserID,
	}

	switch {
	case shoppingCart.ID <= 0:
		// new cart, then insert into transactions table
		queryTrx := `
			INSERT INTO transactions 
				(id_user, total_amount, status) 
			VALUES
				(?, ?, ?)
		`
		// for this query, it should be execute immediately
		_, err := s.client.ExecContext(ctx, queryTrx, shoppingCart.UserID, shoppingCart.GetTotalAmount(), 0)
		if err != nil {
			return nil, fmt.Errorf("unable to create new shopping cart in database due: %w", err)
		}

		// get new transaction / shopping cart ID
		queryTrx = "SELECT id FROM transactions WHERE id_user = ? AND status = 0 LIMIT 1"
		cartRow := dbTx.QueryRowContext(ctx, queryTrx, shoppingCart.UserID)
		var newShoppingCartID int64
		if err = cartRow.Scan(&newShoppingCartID); err != nil {
			return nil, fmt.Errorf("unable to get new shopping cart ID from database due: %w", err)
		}

		// construct query for insert into transaction details table
		completeQueryTrxDetails := s.constructTransactionDetailsQuery(newShoppingCartID, shoppingCart.Details)

		// insert into transaction details table
		_, err = dbTx.ExecContext(ctx, completeQueryTrxDetails)
		if err != nil {
			return nil, fmt.Errorf("unable to insert goods into shopping cart details into database due: %w", err)
		}
	default:
		// construct query for insert into transaction details table
		completeQueryTrxDetails := s.constructTransactionDetailsQuery(shoppingCart.ID, shoppingCart.Details)

		// insert into transaction details first
		_, err = dbTx.ExecContext(ctx, completeQueryTrxDetails)
		if err != nil {
			return nil, fmt.Errorf("unable to insert goods into shopping cart details into database due: %w", err)
		}

		// update total_amount for existing cart in transactions table
		queryTrx := `
			UPDATE transactions
			SET total_amount = (
				SELECT
					SUM(td.total_goods * g.price) AS total_goods_price
				FROM transaction_details td 
				JOIN goods g ON td.id_goods = g.id
			)
			WHERE id = ?
		`
		// update transaction details table
		dbTx.ExecContext(ctx, queryTrx, shoppingCart.ID)
	}

	// commit changes
	err = dbTx.Commit()
	if err != nil {
		return nil, fmt.Errorf("unable to commit add to cart operations in database due: %w", err)
	}

	// fill shopping cart output details
	var latestCart []ShoppingCartRow
	err = s.client.SelectContext(
		ctx,
		&latestCart,
		"SELECT * FROM transactions WHERE id_user = ? AND status = 0 LIMIT 1",
		shoppingCart.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get latest cart details from database due: %w", err)
	}
	simpleCart.ID = latestCart[0].ID
	simpleCart.TotalAmount = latestCart[0].TotalAmount

	return simpleCart, nil
}

func (s storage) constructTransactionDetailsQuery(cartID int64, cartDetails []entity.ShoppingCartDetail) string {
	// prefix query for transaction details
	queryTrxDetails := `
		INSERT INTO transaction_details
			(id_transaction, id_goods, total_goods, created_at)
		VALUES
	`

	trxDetailsValueQueries := []string{}
	for _, goodsDetail := range cartDetails {
		// add multiple values into transaction details query
		vq := fmt.Sprintf(
			`(%d, %d, %d, %d)`,
			cartID,
			goodsDetail.GoodsID,
			goodsDetail.TotalGoods,
			goodsDetail.CreatedAt,
		)
		trxDetailsValueQueries = append(trxDetailsValueQueries, vq)
	}
	completeQueryTrxDetails := []string{
		queryTrxDetails,
		strings.Join(trxDetailsValueQueries, ","),
	}

	return strings.Join(completeQueryTrxDetails, "\n")
}

// CreateTransaction simply update transaction status from `0` to `1`
func (s *storage) CreateTransaction(ctx context.Context, shoppingCart *entity.ShoppingCart) (*entity.Transaction, error) {
	dbTx, err := s.client.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction for add goods to cart query: %w", err)
	}
	// Defer the transaction’s rollback. If the transaction succeeds, it will be committed before the function exits,
	// making the deferred rollback call a no-op.
	// If the transaction fails it won’t be committed, meaning that the rollback will be called as the function exits.
	defer dbTx.Rollback()

	// update transaction status
	queryTrx := `UPDATE transactions SET status = 1 WHERE id = ?`
	_, err = dbTx.ExecContext(ctx, queryTrx, shoppingCart.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to create new transactions into datbaase due: %w", err)
	}

	// get the transaction record to returned it
	var transactionRow struct {
		ID     int64 `db:"id"`
		Status int   `db:"status"`
	}
	queryTrx = `
		SELECT
			id, 
			status 
		FROM transactions
		WHERE id = ? AND status = 1
		LIMIT 1
	`

	trxRow := dbTx.QueryRowContext(ctx, queryTrx, shoppingCart.ID)
	err = trxRow.Scan(
		&transactionRow.ID,
		&transactionRow.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get transaction details due: %w", err)
	}

	// commit changes
	if err = dbTx.Commit(); err != nil {
		return nil, fmt.Errorf("unable to commit creat transaction query in database due: %w", err)
	}

	if transactionRow.Status != 1 {
		return nil, fmt.Errorf("unable to create transaction, shopping cart status not updated properly")
	}

	return &entity.Transaction{
		ID: transactionRow.ID,
	}, nil
}

func (s *storage) TruncateAllData(ctx context.Context) error {
	_, err := s.client.ExecContext(ctx, "TRUNCATE transactions")
	if err != nil {
		return fmt.Errorf("unable to truncate shopping cart / transaction table due: %w", err)
	}
	_, err = s.client.ExecContext(ctx, "TRUNCATE transaction_details")
	if err != nil {
		return fmt.Errorf("unable to truncate shopping cart / transaction details table due: %w", err)
	}
	return nil
}
