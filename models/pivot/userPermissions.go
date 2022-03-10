package pivot

// UserPermissions represents the database model of user permissions relationships
type UserPermissions struct {
	UserID       uint `gorm:"primary_key" json:"user_id"`
	PermissionID uint `gorm:"primary_key" json:"permission_id"`
}

// TableName sets the table name
func (UserPermissions) TableName() string {
	return "user_permissions"
}
