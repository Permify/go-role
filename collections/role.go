package collections

import (
	"github.com/Permify/permify-gorm/helpers"
	"github.com/Permify/permify-gorm/models"
)

// Role provides methods for you to manage array data more easily.
type Role []models.Role

// Origin convert the collection to role array.
// @return []models.Role
func (u Role) Origin() []models.Role {
	return []models.Role(u)
}

// Len returns the number of elements of the array.
// @return int64
func (u Role) Len() (length int64) {
	return int64(len(u))
}

// IDs returns an array of the role array's ids.
// @return []uint
func (u Role) IDs() (IDs []uint) {
	for _, role := range u {
		IDs = append(IDs, role.ID)
	}
	return IDs
}

// Names returns an array of the role array's names.
// @return []string
func (u Role) Names() (names []string) {
	for _, role := range u {
		names = append(names, role.Name)
	}
	return names
}

// GuardNames returns an array of the permission array's guard names.
// @return []string
func (u Role) GuardNames() (guards []string) {
	for _, role := range u {
		guards = append(guards, role.GuardName)
	}
	return guards
}

// Permissions returns uniquely the permissions of the roles in the role collection.
// @return Permission
func (u Role) Permissions() (permissions Permission) {
	var IDs []uint
	for _, a := range u {
		if len(a.Permissions) > 0 {
			for _, prm := range a.Permissions {
				if !helpers.InArray(prm.ID, IDs) {
					permissions = append(permissions, prm)
				}
				IDs = append(IDs, prm.ID)
			}
		}
	}
	return
}
