package auth

import (
	"context"
	"errors"
	apperrors "expense-tracker/internal/shared/errors"

	"gorm.io/gorm"
)

type Repository struct {
	dB *gorm.DB
}

func NewRepository(dB *gorm.DB) *Repository {
	return &Repository{
		dB: dB,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	result := r.dB.WithContext(ctx).Create(user)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return apperrors.Conflict("An account with this email already exists")
		}
		return apperrors.Internal("Failed to create user", result.Error)
	}
	return nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	result := r.dB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("User not found")
		}
		return nil, apperrors.Internal("Failed to find user", result.Error)
	}
	return &user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id uint) (*User, error) {
	var user User
	result := r.dB.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound("User not found")
		}
		return nil, apperrors.Internal("Failed to find user", result.Error)
	}
	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *User) error {
	result := r.dB.WithContext(ctx).Save(user)
	if result.Error != nil {
		return apperrors.Internal("Failed to update user", result.Error)
	}
	return nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, token *RefreshToken) error {
	result := r.dB.WithContext(ctx).Create(token)
	if result.Error != nil {
		return apperrors.Internal("Failed to save refresh token", result.Error)
	}
	return nil
}

func (r *Repository) FindRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	var rt RefreshToken
	result := r.dB.WithContext(ctx).Where("refresh_token = ?", token).First(&rt)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Refresh token not found")
		}
		return nil, apperrors.Internal("Failed to find refresh token", result.Error)
	}
	return &rt, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, token string) error {
	result := r.dB.WithContext(ctx).Where("token = ?", token).Delete(&RefreshToken{})
	if result.Error != nil {
		return apperrors.Internal("Failed to delete refresh token", result.Error)
	}
	return nil
}

func (r *Repository) DeleteAllUserRefreshTokens(ctx context.Context, userID uint) error {
	result := r.dB.WithContext(ctx).Where("user_id = ?", userID).Delete(&User{})
	if result.Error != nil {
		return apperrors.Internal("Failed to delete all user refresh tokens", result.Error)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	return err != nil && (containsString(err.Error(), "duplicate key") || containsString(err.Error(), "unique constraint"))
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
