package options

import (
	`github.com/Permify/permify-gorm/utils`
)

type RoleOption struct {
	WithPermissions bool
	Pagination      *utils.Pagination
}
