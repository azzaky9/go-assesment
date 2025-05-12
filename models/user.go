package models

import "time"

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

type Users struct {
	Id        uint       `gorm:"autoIncrement;primaryKey"`
	Role      *Role      `gorm:"default:user"`
	Products  []Products `gorm:"foreignKey:UserID"`
	Email     string     `gorm:"unique"`
	Password  string     `gorm:"not null"`
	Username  string     `gorm:"unique;not null"`
	FirstName string     `gorm:"not null"`
	LastName  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
