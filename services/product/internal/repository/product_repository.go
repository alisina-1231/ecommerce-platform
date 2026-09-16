package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alisina-1231/ecommerce-platform/services/product/internal/models"
)

var ErrProductNotFound = errors.New("product not found")

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(
	ctx context.Context,
	product *models.Product,
) error {
	const query = `
		INSERT INTO products (
			name,
			description,
			price,
			currency
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Currency,
	).Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	return nil
}

func (r *ProductRepository) GetByID(
	ctx context.Context,
	id int64,
) (*models.Product, error) {
	const query = `
		SELECT
			id,
			name,
			description,
			price,
			currency
		FROM products
		WHERE id = $1;
	`

	var product models.Product

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Currency,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}

		return nil, fmt.Errorf("get product: %w", err)
	}

	return &product, nil
}

func (r *ProductRepository) List(
	ctx context.Context,
) ([]models.Product, error) {
	const query = `
		SELECT
			id,
			name,
			description,
			price,
			currency
		FROM products
		ORDER BY id;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		var product models.Product

		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Currency,
		); err != nil {
			return nil, fmt.Errorf(
				"scan product: %w",
				err,
			)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate products: %w",
			err,
		)
	}

	return products, nil
}

func (r *ProductRepository) Update(
	ctx context.Context,
	product *models.Product,
) error {
	const query = `
		UPDATE products
		SET
			name = $1,
			description = $2,
			price = $3,
			currency = $4,
			updated_at = NOW()
		WHERE id = $5;
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Currency,
		product.ID,
	)

	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"check update result: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return ErrProductNotFound
	}

	return nil
}

func (r *ProductRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	const query = `
		DELETE FROM products
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"check delete result: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return ErrProductNotFound
	}

	return nil
}
