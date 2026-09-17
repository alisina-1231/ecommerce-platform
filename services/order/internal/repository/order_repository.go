package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alisina-1231/ecommerce-platform/services/order/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrOrderNotFound = errors.New("order not found")

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(
	ctx context.Context,
	request models.CreateOrderRequest,
) (*models.Order, error) {
	query := `
		INSERT INTO orders (
			user_id,
			total_amount,
			currency,
			status,
			shipping_address
		)
		VALUES ($1, $2, $3, 'PENDING', $4)
		RETURNING
			id,
			user_id,
			total_amount,
			currency,
			status,
			shipping_address,
			created_at,
			updated_at
	`

	var order models.Order

	err := r.db.QueryRow(
		ctx,
		query,
		request.UserID,
		request.TotalAmount,
		request.Currency,
		request.ShippingAddress,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.Currency,
		&order.Status,
		&order.ShippingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) GetByID(
	ctx context.Context,
	id int64,
) (*models.Order, error) {
	query := `
		SELECT
			id,
			user_id,
			total_amount,
			currency,
			status,
			shipping_address,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
	`

	var order models.Order

	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.Currency,
		&order.Status,
		&order.ShippingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrOrderNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) GetByUserID(
	ctx context.Context,
	userID string,
) ([]models.Order, error) {
	query := `
		SELECT
			id,
			user_id,
			total_amount,
			currency,
			status,
			shipping_address,
			created_at,
			updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("get user orders: %w", err)
	}

	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {
		var order models.Order

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.Currency,
			&order.Status,
			&order.ShippingAddress,
			&order.CreatedAt,
			&order.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateStatus(
	ctx context.Context,
	id int64,
	status string,
) (*models.Order, error) {
	query := `
		UPDATE orders
		SET
			status = $2,
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id,
			user_id,
			total_amount,
			currency,
			status,
			shipping_address,
			created_at,
			updated_at
	`

	var order models.Order

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		status,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.Currency,
		&order.Status,
		&order.ShippingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrOrderNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("update order status: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	result, err := r.db.Exec(
		ctx,
		"DELETE FROM orders WHERE id = $1",
		id,
	)

	if err != nil {
		return fmt.Errorf("delete order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	return nil
}
