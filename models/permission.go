package models

import (
	"time"
)

// Permission represents the database model of permissions
type Permission struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Name        string `gorm:"size:255;not null" json:"name"`
	GuardName   string `gorm:"size:255;not null;index" json:"guard_name"`
	Description string `gorm:"size:255" json:"description"`

	// Time
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sets the table name
func (Permission) TableName() string {
	return "permissions"
}
