package categories

import (
	"context"
	"errors"
	apperrors "expense-tracker/internal/shared/errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, c *Category) error {
	result := r.db.WithContext(ctx).Create(c)
	if result.Error != nil {
		return apperrors.Internal("failed to create category", result.Error)
	}
	return nil
}

func (r *Repository) FindAllByUserID(ctx context.Context, userID uint) ([]Category, error) {
	var categories []Category
	result := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).Order("name ASC").Find(&categories)
	if result.Error != nil {
		return nil, apperrors.Internal("failed to find all categories", result.Error)
	}
	return categories, nil
}

func (r *Repository) FindByIDAndUserID(ctx context.Context, id, userID uint) (*Category, error) {
	var category Category
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Category not found")
		}
		return nil, apperrors.Internal("failed to fetch category", result.Error)
	}
	return &category, nil
}

// FindByID Used by internal services like budgets, expenses
func (r *Repository) FindByID(ctx context.Context, id uint) (*Category, error) {
	var category Category
	result := r.db.WithContext(ctx).First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Category")
		}
		return nil, apperrors.Internal("failed to fetch category", result.Error)
	}
	return &category, nil
}

func (r *Repository) Update(ctx context.Context, c *Category) error {
	result := r.db.WithContext(ctx).Save(c)
	if result.Error != nil {
		return apperrors.Internal("failed to update category", result.Error)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id, userID uint) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ? AND is_default = false", id, userID).Delete(&Category{})
	if result.Error != nil {
		return apperrors.Internal("failed to delete category", result.Error)
	}

	if result.RowsAffected == 0 {
		return apperrors.BadRequest("Category not found or cannot delete a default category")
	}

	return nil
}

func (r *Repository) SeedDefaultCategories(ctx context.Context, userID uint) error {
	defaults := defaultCategories(userID)
	result := r.db.WithContext(ctx).CreateInBatches(defaults, len(defaults))
	if result.Error != nil {
		return apperrors.Internal("failed to create default categories", result.Error)
	}
	return nil
}

func (r *Repository) UserHasCategory(ctx context.Context, userID uint, categoryID uint) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&Category{}).Where("id = ? AND user_id = ?", categoryID, userID).Count(&count)
	if result.Error != nil {
		return false, apperrors.Internal("Failed to verify category ownership", result.Error)
	}
	return count > 0, nil
}
