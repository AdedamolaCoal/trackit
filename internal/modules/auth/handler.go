package auth

import (
	apperrors "expense-tracker/internal/shared/errors"
	"expense-tracker/internal/shared/middleware"
	"expense-tracker/pkg/response"
	"expense-tracker/pkg/validator"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register REGISTER
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if errs := validator.DecodeAndValidate(r, &req); errs != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeBadRequest, "Validation Failed", errs)
		return
	}

	auth, err := h.service.Register(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "Account created successfully", auth)
}

// Login LOGIN
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if errs := validator.DecodeAndValidate(r, &req); errs != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeBadRequest, "Validation Failed", errs)
		return
	}

	auth, err := h.service.Login(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Login successful", auth)
}

// Refresh REFRESH
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if errs := validator.DecodeAndValidate(r, &req); errs != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeBadRequest, "Validation Failed", errs)
		return
	}

	auth, err := h.service.RefreshToken(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Refresh Token successfully", auth)
}

// Logout LOGOUT
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if errs := validator.DecodeAndValidate(r, &req); errs != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeBadRequest, "Validation Failed", errs)
		return
	}

	if err := h.service.Logout(r.Context(), req.RefreshToken); err != nil {
		handleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Logout successfully", nil)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "unauthorized")
		return
	}

	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "profile retrieved successfully", ToUserResponse(user))
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "unauthorized")
		return
	}

	var req UpdateProfileRequest
	if errs := validator.DecodeAndValidate(r, &req); errs != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeValidation, "validation failed", errs)
		return
	}

	user, err := h.service.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "profile updated successfully", ToUserResponse(user))
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "unauthorized")
		return
	}

	var req ChangePasswordRequest
	if errs := validator.DecodeAndValidate(r, &req); errs != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeValidation, "validation failed", errs)
		return
	}

	if err := h.service.ChangePassword(r.Context(), userID, req); err != nil {
		handleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "password changed successfully", nil)
}

// handleError is shared across all handlers in this package —
// converts AppError to the right response, handles unknown errors uniformly
func handleError(w http.ResponseWriter, err error) {
	var appErr *apperrors.AppError
	if apperrors.As(err, &appErr) {
		response.Error(w, appErr.HttpStatus, appErr.Code, appErr.Message)
		return
	}
	response.Error(w, http.StatusInternalServerError, response.ErrCodeInternal, "something went wrong")
}
