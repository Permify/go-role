package models_pivot

type UserPermissions struct {
	UserID       uint `gorm:"primary_key" json:"user_id"`
	PermissionID uint `gorm:"primary_key" json:"permission_id"`
}

// TableName /**
func (UserPermissions) TableName() string {
	return "user_permissions"
}
