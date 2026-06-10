package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"expense-tracker/internal/config"
	apperrors "expense-tracker/internal/shared/errors"
	"expense-tracker/internal/shared/middleware"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	cfg  *config.Config
	repo *Repository
}

func NewService(cfg *config.Config, repo *Repository) *Service {
	return &Service{cfg: cfg, repo: repo}
}

// Register REGISTER
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Internal("Failed to hash password", err)
	}
	currency := req.Currency

	if currency == "" {
		currency = "NGN"
	}

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Currency: currency,
		Password: string(hashed),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

// Login LOGIN
func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.InvalidCredentials()
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.InvalidCredentials()
	}

	return s.generateAuthResponse(ctx, user)
}

// RefreshToken REFRESH TOKEN
func (s *Service) RefreshToken(ctx context.Context, req RefreshRequest) (*AuthResponse, error) {
	stored, err := s.repo.FindRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if stored.IsTokenExpired() {
		_ = s.repo.DeleteRefreshToken(ctx, req.RefreshToken)
		return nil, apperrors.TokenExpired()
	}

	user, err := s.repo.FindUserByID(ctx, stored.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteRefreshToken(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

// Logout LOGOUT
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.DeleteRefreshToken(ctx, refreshToken)
}

// GetProfile GET PROFILE
func (s *Service) GetProfile(ctx context.Context, userID uint) (*User, error) {
	return s.repo.FindUserByID(ctx, userID)
}

// UpdateProfile UPDATE PROFILE
func (s *Service) UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) (*User, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Currency != "" {
		user.Currency = req.Currency
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword CHANGE PASSWORD
func (s *Service) ChangePassword(ctx context.Context, userID uint, req ChangePasswordRequest) error {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return apperrors.BadRequest("Current Password is incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.Internal("Failed to hash password", err)
	}

	user.Password = string(hashed)
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return err
	}
	return s.repo.DeleteAllUserRefreshTokens(ctx, userID)
}

// Helpers HELPERS
func (s *Service) generateAuthResponse(ctx context.Context, user *User) (*AuthResponse, error) {
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateAndStoreRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         ToUserResponse(user),
	}, nil
}

func (s *Service) generateAccessToken(user *User) (string, error) {
	claims := middleware.Claims{
		UserID: strconv.Itoa(int(user.ID)),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.Jwt.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "access",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.Jwt.Secret))

	if err != nil {
		return "", apperrors.Internal("Failed to sign token", err)
	}
	return signed, nil
}

func (s *Service) generateAndStoreRefreshToken(ctx context.Context, userID uint) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", apperrors.Internal("Failed to generate refresh token", err)
	}
	tokenStr := hex.EncodeToString(raw)
	rt := &RefreshToken{
		UserID:    userID,
		Token:     tokenStr,
		ExpiredAt: time.Now().Add(s.cfg.Jwt.RefreshTTL),
	}

	if err := s.repo.SaveRefreshToken(ctx, rt); err != nil {
		return "", err
	}
	return tokenStr, nil
}
