package categories

import (
	"time"

	"gorm.io/gorm"
)

type CategoryType string
type Category struct {
	ID        uint           `gorm:"primary_key;autoIncrement" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Name      string         `gorm:"type:varchar(100);not null;index" json:"name"`
	Icon      string         `gorm:"type:varchar(50)" json:"icon"`
	Color     string         `gorm:"type:varchar(7)" json:"color"`
	Type      CategoryType   `gorm:"type:varchar(100);not null;default:'expense'" json:"type"`
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	CategoryTypeExpense CategoryType = "expense"
	CategoryTypeIncome  CategoryType = "income"
)
