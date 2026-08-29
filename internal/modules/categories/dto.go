package categories

type CreateCategoryRequest struct {
	Name  string       `json:"name" validate:"required,min=1,max=100"`
	Icon  string       `json:"icon" validate:"omitempty,max=50"`
	Color string       `json:"color" validate:"omitempty,len=7"`
	Type  CategoryType `json:"type" validate:"required,oneof=expense income"`
}

type UpdateCategoryRequest struct {
	Name  string       `json:"name" validate:"required,min=1,max=100"`
	Icon  string       `json:"icon" validate:"omitempty,max=50"`
	Color string       `json:"color" validate:"omitempty,len=7"`
	Type  CategoryType `json:"type" validate:"required,oneof=expense income"`
}

type CategoryResponse struct {
	ID        uint         `json:"id"`
	Name      string       `json:"name"`
	Icon      string       `json:"icon"`
	Color     string       `json:"color"`
	Type      CategoryType `json:"type"`
	IsDefault bool         `json:"is_default"`
}

func ToCategoryResponse(c *Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Icon:      c.Icon,
		Color:     c.Color,
		Type:      c.Type,
		IsDefault: c.IsDefault,
	}
}

func ToCategoryResponseList(categories []Category) []CategoryResponse {
	result := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		result[i] = ToCategoryResponse(&category)
	}
	return result
}
