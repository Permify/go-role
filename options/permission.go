package options

import (
	`github.com/Permify/permify-gorm/utils`
)

// PermissionOption represents options when fetching permissions.
type PermissionOption struct {
	Pagination *utils.Pagination
}