package services

import (
	"context"
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
