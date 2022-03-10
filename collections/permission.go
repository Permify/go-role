package collections

import (
	"github.com/Permify/permify-gorm/models"
)

// Permission provides methods for you to manage array data more easily.
type Permission []models.Permission

// Origin convert the collection to permission array.
// @return []models.Permission
func (u Permission) Origin() []models.Permission {
	return []models.Permission(u)
}

// Len returns the number of elements of the array.
// @return int64
func (u Permission) Len() (length int64) {
	return int64(len(u))
}

// IDs returns an array of the permission array's ids.
// @return []uint
func (u Permission) IDs() (IDs []uint) {
	for _, permission := range u {
		IDs = append(IDs, permission.ID)
	}
	return IDs
}

// Names returns an array of the permission array's names.
// @return []string
func (u Permission) Names() (names []string) {
	for _, permission := range u {
		names = append(names, permission.Name)
	}
	return names
}

// GuardNames returns an array of the permission array's guard names.
// @return []string
func (u Permission) GuardNames() (guards []string) {
	for _, permission := range u {
		guards = append(guards, permission.GuardName)
	}
	return guards
}
