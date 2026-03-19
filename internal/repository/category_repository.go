package repository

import (
	"context"
	"database/sql"
	"money-tracker/internal/models"

	"github.com/google/uuid"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		DB: db,
	}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, c *models.Category) error {
	query := `INSERT INTO categories (id, user_id, allocation_id, name, type, icon, color) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`

	err := r.DB.QueryRowContext(ctx, query, c.UserID, c.AllocationID, c.Name, c.Type, c.Icon, c.Color).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)

	return err
}

func (r *CategoryRepository) GetCategories(ctx context.Context, userID uuid.UUID) ([]*models.Category, error) {
	query := `SELECT id, allocation_id, name, type, icon, color FROM categories WHERE user_id = $1`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := rows.Close()
		if err == nil {
			err = closeErr
		}
	}()

	var categories []*models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.AllocationID, &category.Name, &category.Type, &category.Icon, &category.Color)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if categories == nil {
		categories = []*models.Category{}
	}

	return categories, nil
}
