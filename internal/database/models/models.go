package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null" binding:"required,min=3,max=100"`
	Content   string    `json:"content" gorm:"not null" binding:"required,min=10"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Custom validator for Gin
var validate = validator.New()

func (p *Post) Validate() error {
	return validate.Struct(p)
}
