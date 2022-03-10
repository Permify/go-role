package options

import (
	"github.com/Permify/permify-gorm/utils"
)

// RoleOption represents options when fetching roles.
type RoleOption struct {
	WithPermissions bool
	Pagination      *utils.Pagination
}
