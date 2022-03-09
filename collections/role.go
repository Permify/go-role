package collections

import (
	`github.com/Permify/permify-gorm/helpers`
	`github.com/Permify/permify-gorm/models`
)

// Role */
type Role []models.Role

// Permissions */
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

// Origin */
func (u Role) Origin() []models.Role {
	return []models.Role(u)
}

// Len */
func (u Role) Len() (length int64) {
	return int64(len(u))
}

// IDs */
func (u Role) IDs() (IDs []uint) {
	for _, role := range u {
		IDs = append(IDs, role.ID)
	}
	return IDs
}

// Names */
func (u Role) Names() (names []string) {
	for _, role := range u {
		names = append(names, role.Name)
	}
	return names
}

// GuardNames */
func (u Role) GuardNames() (guards []string) {
	for _, role := range u {
		guards = append(guards, role.GuardName)
	}
	return guards
}
