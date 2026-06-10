package auth

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Currency  string         `gorm:"type:varchar(10);default:'NGN'" json:"currency"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(512);uniqueIndex;not null"`
	ExpiredAt time.Time `gorm:"expiresAt;not null"`
	CreatedAt time.Time
	User      User `gorm:"foreignKey;UserID"`
}

//func (RefreshToken) IsExpired() bool {
//	return time.Now().After(RefreshToken{}.ExpiredAt)
//}

func (t *RefreshToken) IsTokenExpired() bool {
	return time.Now().After(t.ExpiredAt)
}
