package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type PaginatedData struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func JSON(w http.ResponseWriter, status int, payload Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		return
	}
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	JSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, status int, code, message string, details ...interface{}) {
	var d interface{}
	if len(details) > 0 {
		d = details[0]
	}

	JSON(w, status, Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: d,
		},
	})
}

func Paginated(w http.ResponseWriter, status int, page, pageSize int, message string, items interface{}, total int64) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize == 0 {
		totalPages++
	}

	Success(w, status, message, PaginatedData{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}
