package middleware

import (
	"context"
	"net/http"
	"strings"

	"expense-tracker/internal/config"
	"expense-tracker/internal/shared/errors"
	"expense-tracker/pkg/response"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

type AuthMiddleware struct {
	cfg *config.Config
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractToken(r)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, err.Error())
			return
		}

		claims, appErr := validateToken(token, m.cfg.Jwt.Secret)
		if appErr != nil {
			response.Error(w, appErr.HttpStatus, appErr.Code, appErr.Message)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.Unauthorized("missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", errors.Unauthorized("Authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

func validateToken(tokenStr, secret string) (*Claims, *errors.AppError) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.InvalidToken("Unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, errors.TokenExpired()
		}
		return nil, errors.InvalidToken("Invalid or malformed token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.InvalidToken("Could not parse token claims")
	}
	return claims, nil
}

func GetUserID(r *http.Request) (uint, bool) {
	id, ok := r.Context().Value(UserIDKey).(uint)
	return id, ok
}
