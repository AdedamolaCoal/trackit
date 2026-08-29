package categories

import (
	"context"
	apperrors "expense-tracker/internal/shared/errors"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateCategory(ctx context.Context, userID uint, req CreateCategoryRequest) (*Category, error) {
	category := &Category{
		UserID: userID,
		Name:   req.Name,
		Icon:   req.Icon,
		Color:  req.Color,
		Type:   req.Type,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id, userID uint, req UpdateCategoryRequest) (*Category, error) {
	category, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if category.IsDefault {
		return nil, apperrors.Forbidden("default categories cannot be modified")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.Color != "" {
		category.Color = req.Color
	}
	if req.Type != "" {
		category.Type = req.Type
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Service) GetAllCategories(ctx context.Context, userID uint) ([]Category, error) {
	return s.repo.FindAllByUserID(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id, userID uint) (*Category, error) {
	return s.repo.FindByIDAndUserID(ctx, id, userID)
}

func (s *Service) DeleteCategory(ctx context.Context, userID, id uint) error {
	return s.repo.Delete(ctx, id, userID)
}

// SeedDefaultCategories Seed by Default is called once after registration
func (s *Service) SeedDefaultCategories(ctx context.Context, userID uint) error {
	return s.repo.SeedDefaultCategories(ctx, userID)
}

// ValidateOwnership is used by expenses and budgets to verify
// a category belongs to the user before allowing it to be assigned
func (s *Service) ValidateOwnership(ctx context.Context, categoryID, userID uint) error {
	owns, err := s.repo.UserHasCategory(ctx, categoryID, userID)
	if err != nil {
		return err
	}
	if !owns {
		return apperrors.Forbidden("category does not belong to this user")
	}
	return nil
}
