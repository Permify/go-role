package collections

import (
	`github.com/Permify/permify-gorm/models`
)

// Permission */
type Permission []models.Permission

// Origin */
func (u Permission) Origin() []models.Permission {
	return []models.Permission(u)
}

// Len */
func (u Permission) Len() (length int64) {
	return int64(len(u))
}

// IDs */
func (u Permission) IDs() (IDs []uint) {
	for _, permission := range u {
		IDs = append(IDs, permission.ID)
	}
	return IDs
}

// Names */
func (u Permission) Names() (names []string) {
	for _, permission := range u {
		names = append(names, permission.Name)
	}
	return names
}

// GuardNames */
func (u Permission) GuardNames() (guards []string) {
	for _, permission := range u {
		guards = append(guards, permission.GuardName)
	}
	return guards
}
