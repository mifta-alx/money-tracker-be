package services

import (
	"context"
	"database/sql"
	"errors"
	"money-tracker/internal/models"
	"money-tracker/internal/repository"

	"github.com/google/uuid"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(r *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *models.Category) (*models.Category, error) {
	if req.UserID == uuid.Nil || req.Name == "" || req.Icon == "" {
		return nil, ErrMissingRequiredFields
	}

	err := s.repo.CreateCategory(ctx, req)
	if err != nil {
		return nil, ErrInternal
	}

	return req, nil
}

func (s *CategoryService) GetCategories(ctx context.Context, userID uuid.UUID) ([]*models.Category, error) {
	categories, err := s.repo.GetCategories(ctx, userID)
	if err != nil {
		return nil, ErrInternal
	}
	return categories, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, categoryID, userID uuid.UUID) (*models.Category, error) {
	category, err := s.repo.GetCategory(ctx, categoryID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
	}
	return category, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, req *models.Category) (*models.Category, error) {
	err := s.repo.UpdateCategory(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, ErrInternal
	}
	return req, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, categoryID, userID uuid.UUID) error {
	err := s.repo.DeleteCategory(ctx, categoryID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return ErrInternal
	}
	return nil
}
