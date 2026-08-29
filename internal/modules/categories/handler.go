package categories

import (
	apperrors "expense-tracker/internal/shared/errors"
	"expense-tracker/internal/shared/middleware"
	"expense-tracker/pkg/response"
	"expense-tracker/pkg/validator"
	"net/http"
	"strconv"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	var req CreateCategoryRequest
	if errs := validator.DecodeAndValidate(r, &req); errs != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeValidation, "validation failed", errs)
		return
	}

	category, err := h.s.CreateCategory(r.Context(), userID, req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "Category created successfully", category)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "unauthorized")
		return
	}

	categories, err := h.s.GetAllCategories(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "categories retrieved successfully", ToCategoryResponseList(categories))
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "unauthorized")
		return
	}

	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeBadRequest, "invalid category id")
		return
	}

	category, svcErr := h.s.GetByID(r.Context(), id, userID)
	if svcErr != nil {
		handleError(w, svcErr)
		return
	}

	response.Success(w, http.StatusOK, "category retrieved successfully", ToCategoryResponse(category))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "unauthorized")
		return
	}

	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeBadRequest, "invalid category id")
		return
	}

	var req UpdateCategoryRequest
	if errs := validator.DecodeAndValidate(r, &req); errs != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeValidation, "validation failed", errs)
		return
	}

	category, svcErr := h.s.UpdateCategory(r.Context(), id, userID, req)
	if svcErr != nil {
		handleError(w, svcErr)
		return
	}

	response.Success(w, http.StatusOK, "category updated successfully", ToCategoryResponse(category))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, response.ErrCodeUnauthorized, "unauthorized")
		return
	}

	id, err := parseID(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, response.ErrCodeBadRequest, "invalid category id")
		return
	}

	if svcErr := h.s.DeleteCategory(r.Context(), id, userID); svcErr != nil {
		handleError(w, svcErr)
		return
	}

	response.Success(w, http.StatusOK, "category deleted successfully", nil)
}

func handleError(w http.ResponseWriter, err error) {
	var appErr *apperrors.AppError
	if apperrors.As(err, &appErr) {
		response.Error(w, appErr.HttpStatus, appErr.Code, appErr.Message)
		return
	}
	response.Error(w, http.StatusInternalServerError, response.ErrCodeInternal, "something went wrong")
}

func parseID(r *http.Request, key string) (uint, error) {
	val := r.PathValue(key)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
