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
	query := `INSERT INTO categories (user_id, allocation_id, name, type, icon, color, target_amount) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	return r.DB.QueryRowContext(ctx, query, c.UserID, c.AllocationID, c.Name, c.Type, c.Icon, c.Color, c.TargetAmount).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, c *models.Category) error {
	query := `UPDATE categories SET allocation_id = $1, name = $2, type = $3, icon = $4, color = $5, target_amount = $6, updated_at = NOW() WHERE id = $7 AND user_id = $8 RETURNING id, created_at, updated_at`
	return r.DB.QueryRowContext(ctx, query, c.AllocationID, c.Name, c.Type, c.Icon, c.Color, c.TargetAmount, c.ID, c.UserID).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`
	result, err := r.DB.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (r *CategoryRepository) GetCategories(ctx context.Context, userID uuid.UUID) ([]*models.Category, error) {
	query := `SELECT id, allocation_id, name, type, icon, color, target_amount, created_at, updated_at FROM categories WHERE user_id = $1 ORDER BY created_at`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var categories []*models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.AllocationID, &category.Name, &category.Type, &category.Icon, &category.Color, &category.TargetAmount, &category.CreatedAt, &category.UpdatedAt)
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

func (r *CategoryRepository) GetCategory(ctx context.Context, id, userID uuid.UUID) (*models.Category, error) {
	query := `SELECT id, allocation_id, name, type, icon, color, target_amount, created_at, updated_at FROM categories WHERE id = $1 AND user_id = $2`
	var category models.Category
	err := r.DB.QueryRowContext(ctx, query, id, userID).Scan(&category.ID, &category.AllocationID, &category.Name, &category.Type, &category.Icon, &category.Color, &category.TargetAmount, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
