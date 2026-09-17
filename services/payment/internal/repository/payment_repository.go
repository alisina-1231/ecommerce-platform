package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alisina-1231/ecommerce-platform/services/payment/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPaymentNotFound = errors.New("payment not found")

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) Create(
	ctx context.Context,
	request models.CreatePaymentRequest,
) (*models.Payment, error) {
	query := `
		INSERT INTO payments (
			order_id,
			user_id,
			amount,
			currency,
			payment_method,
			status
		)
		VALUES ($1, $2, $3, $4, $5, 'PENDING')
		RETURNING
			id,
			order_id,
			user_id,
			amount,
			currency,
			payment_method,
			status,
			transaction_id,
			failure_reason,
			created_at,
			updated_at
	`

	var payment models.Payment

	err := r.db.QueryRow(
		ctx,
		query,
		request.OrderID,
		request.UserID,
		request.Amount,
		request.Currency,
		request.PaymentMethod,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.PaymentMethod,
		&payment.Status,
		&payment.TransactionID,
		&payment.FailureReason,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	return &payment, nil
}

func (r *PaymentRepository) GetByID(
	ctx context.Context,
	id int64,
) (*models.Payment, error) {
	query := `
		SELECT
			id,
			order_id,
			user_id,
			amount,
			currency,
			payment_method,
			status,
			transaction_id,
			failure_reason,
			created_at,
			updated_at
		FROM payments
		WHERE id = $1
	`

	var payment models.Payment

	err := r.db.QueryRow(ctx, query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.PaymentMethod,
		&payment.Status,
		&payment.TransactionID,
		&payment.FailureReason,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPaymentNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}

	return &payment, nil
}

func (r *PaymentRepository) GetByUserID(
	ctx context.Context,
	userID string,
) ([]models.Payment, error) {
	query := `
		SELECT
			id,
			order_id,
			user_id,
			amount,
			currency,
			payment_method,
			status,
			transaction_id,
			failure_reason,
			created_at,
			updated_at
		FROM payments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user payments: %w", err)
	}
	defer rows.Close()

	payments := make([]models.Payment, 0)

	for rows.Next() {
		var payment models.Payment

		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.UserID,
			&payment.Amount,
			&payment.Currency,
			&payment.PaymentMethod,
			&payment.Status,
			&payment.TransactionID,
			&payment.FailureReason,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan payment: %w", err)
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate payments: %w", err)
	}

	return payments, nil
}

func (r *PaymentRepository) UpdateStatus(
	ctx context.Context,
	id int64,
	request models.UpdatePaymentStatusRequest,
) (*models.Payment, error) {
	query := `
		UPDATE payments
		SET
			status = $2,
			transaction_id = COALESCE($3, transaction_id),
			failure_reason = COALESCE($4, failure_reason),
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id,
			order_id,
			user_id,
			amount,
			currency,
			payment_method,
			status,
			transaction_id,
			failure_reason,
			created_at,
			updated_at
	`

	var payment models.Payment

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		request.Status,
		request.TransactionID,
		request.FailureReason,
	).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.PaymentMethod,
		&payment.Status,
		&payment.TransactionID,
		&payment.FailureReason,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPaymentNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("update payment status: %w", err)
	}

	return &payment, nil
}

func (r *PaymentRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	commandTag, err := r.db.Exec(
		ctx,
		"DELETE FROM payments WHERE id = $1",
		id,
	)

	if err != nil {
		return fmt.Errorf("delete payment: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrPaymentNotFound
	}

	return nil
}
