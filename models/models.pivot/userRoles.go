package models_pivot

type UserRoles struct {
	UserID uint `gorm:"primary_key" json:"user_id"`
	RoleID uint `gorm:"primary_key" json:"role_id"`
}

// TableName /**
func (UserRoles) TableName() string {
	return "user_roles"
}