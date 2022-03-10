package repositories

import (
	`gorm.io/gorm`
	`gorm.io/gorm/clause`

	`github.com/Permify/permify-gorm/collections`
	`github.com/Permify/permify-gorm/models`
	models_pivot `github.com/Permify/permify-gorm/models/models.pivot`
)

type IUserRepository interface {

	// actions

	AddPermissions(userID uint, permissions collections.Permission) (err error)
	ReplacePermissions(userID uint, permissions collections.Permission) (err error)
	RemovePermissions(userID uint, permissions collections.Permission) (err error)
	ClearPermissions(userID uint) (err error)

	AddRoles(userID uint, roles collections.Role) (err error)
	ReplaceRoles(userID uint, roles collections.Role) (err error)
	RemoveRoles(userID uint, roles collections.Role) (err error)
	ClearRoles(userID uint) (err error)

	// controls

	HasRole(userID uint, role models.Role) (b bool, err error)
	HasAllRoles(userID uint, roles collections.Role) (b bool, err error)
	HasAnyRoles(userID uint, roles collections.Role) (b bool, err error)

	HasDirectPermission(userID uint, permission models.Permission) (b bool, err error)
	HasAllDirectPermissions(userID uint, permissions collections.Permission) (b bool, err error)
	HasAnyDirectPermissions(userID uint, permissions collections.Permission) (b bool, err error)
}

type UserRepository struct {
	Database *gorm.DB
}

// ACTIONS

// AddPermissions
// Add direct permissions to user.
// @param uint
// @param collections.Permission
// @return error
func (repository *UserRepository) AddPermissions(userID uint, permissions collections.Permission) error {
	return repository.Database.Transaction(func(tx *gorm.DB) error {

		var ups []models_pivot.UserPermissions

		for _, p := range permissions.Origin() {
			ups = append(ups, models_pivot.UserPermissions{
				UserID:       userID,
				PermissionID: p.ID,
			})
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&ups).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

// ReplacePermissions
// Replace direct permissions of user.
// @param uint
// @param collections.Permission
// @return error
func (repository *UserRepository) ReplacePermissions(userID uint, permissions collections.Permission) error {
	return repository.Database.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("user_permissions.user_id = ?", userID).Delete(&models_pivot.UserPermissions{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		var ups []models_pivot.UserPermissions

		for _, p := range permissions.Origin() {
			ups = append(ups, models_pivot.UserPermissions{
				UserID:       userID,
				PermissionID: p.ID,
			})
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&ups).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

// RemovePermissions
// Remove direct permissions from user.
// @param uint
// @param collections.Permission
// @return error
func (repository *UserRepository) RemovePermissions(userID uint, permissions collections.Permission) error {
	return repository.Database.Transaction(func(tx *gorm.DB) error {

		var ups []models_pivot.UserPermissions

		for _, p := range permissions.Origin() {
			ups = append(ups, models_pivot.UserPermissions{
				UserID:       userID,
				PermissionID: p.ID,
			})
		}

		if err := tx.Delete(&ups).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

// ClearPermissions
// Remove all direct permissions from user.
// @param uint
// @return error
func (repository *UserRepository) ClearPermissions(userID uint) (err error) {
	return repository.Database.Where("user_permissions.user_id = ?", userID).Delete(&models_pivot.UserPermissions{}).Error
}

// AddRoles
// Add roles to user.
// @param uint
// @param collections.Role
// @return error
func (repository *UserRepository) AddRoles(userID uint, roles collections.Role) error {
	return repository.Database.Transaction(func(tx *gorm.DB) error {
		var ur []models_pivot.UserRoles
		for _, r := range roles.Origin() {
			ur = append(ur, models_pivot.UserRoles{
				UserID: userID,
				RoleID: r.ID,
			})
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&ur).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

// ReplaceRoles
// Replace roles of user.
// @param uint
// @param collections.Role
// @return error
func (repository *UserRepository) ReplaceRoles(userID uint, roles collections.Role) error {
	return repository.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_roles.user_id = ?", userID).Delete(&models_pivot.UserRoles{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		var ur []models_pivot.UserRoles
		for _, r := range roles.Origin() {
			ur = append(ur, models_pivot.UserRoles{
				UserID: userID,
				RoleID: r.ID,
			})
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&ur).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

// RemoveRoles
// Remove roles from user.
// @param uint
// @param collections.Role
// @return error
func (repository *UserRepository) RemoveRoles(userID uint, roles collections.Role) error {
	return repository.Database.Transaction(func(tx *gorm.DB) error {
		var ur []models_pivot.UserRoles
		for _, r := range roles.Origin() {
			ur = append(ur, models_pivot.UserRoles{
				UserID: userID,
				RoleID: r.ID,
			})
		}
		if err := tx.Delete(&ur).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

// ClearRoles
// Remove all roles from user.
// @param uint
// @return error
func (repository *UserRepository) ClearRoles(userID uint) (err error) {
	return repository.Database.Where("user_roles.user_id = ?", userID).Delete(&models_pivot.UserRoles{}).Error
}

// CONTROLS

// HasRole
// Does the user have the given role?
// @param uint
// @param models.Role
// @return bool, error
func (repository *UserRepository) HasRole(userID uint, role models.Role) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Where("user_roles.role_id = ?", role.ID).Count(&count).Error
	return count > 0, err
}

// HasAllRoles
// Does the user have all the given roles?
// @param uint
// @param collections.Role
// @return bool, error
func (repository *UserRepository) HasAllRoles(userID uint, roles collections.Role) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Where("user_roles.role_id IN (?)", roles.IDs()).Count(&count).Error
	return roles.Len() == count, err
}

// HasAnyRoles
// Does the user have any of the given roles?
// @param uint
// @param collections.Role
// @return bool, error
func (repository *UserRepository) HasAnyRoles(userID uint, roles collections.Role) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Where("user_roles.role_id IN (?)", roles.IDs()).Count(&count).Error
	return count > 0, err
}

// HasDirectPermission
// Does the user have the given permission? (not including the permissions of the roles)
// @param uint
// @param collections.Permission
// @return bool, error
func (repository *UserRepository) HasDirectPermission(userID uint, permission models.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Where("user_permissions.permission_id = ?", permission.ID).Count(&count).Error
	return count > 0, err
}

// HasAllDirectPermissions
// Does the user have all the given permissions? (not including the permissions of the roles)
// @param uint
// @param collections.Permission
// @return bool, error
func (repository *UserRepository) HasAllDirectPermissions(userID uint, permissions collections.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Where("user_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return permissions.Len() == count, err
}

// HasAnyDirectPermissions
// Does the user have any of the given permissions? (not including the permissions of the roles)
// @param uint
// @param collections.Permission
// @return bool, error
func (repository *UserRepository) HasAnyDirectPermissions(userID uint, permissions collections.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Where("user_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return count > 0, err
}
