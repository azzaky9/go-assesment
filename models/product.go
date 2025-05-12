package models

import "time"

type Products struct {
	Id        uint `gorm:"autoIncrement"`
	Name      string
	Quantity  uint
	Price     float64
	UserID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
