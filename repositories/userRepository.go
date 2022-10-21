package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Permify/go-role/collections"
	"github.com/Permify/go-role/models"
	"github.com/Permify/go-role/models/pivot"
)

// IUserRepository its data access layer abstraction of user.
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

// UserRepository its data access layer of user.
type UserRepository struct {
	Database *gorm.DB
}

// ACTIONS

// AddPermissions add direct permissions to user.
// @param uint
// @param collections.Permission
// @return error
func (repository *UserRepository) AddPermissions(userID uint, permissions collections.Permission) error {
	var userPermissions []pivot.UserPermissions
	for _, permission := range permissions.Origin() {
		userPermissions = append(userPermissions, pivot.UserPermissions{
			UserID:       userID,
			PermissionID: permission.ID,
		})
	}
	return repository.Database.Clauses(clause.OnConflict{DoNothing: true}).Create(&userPermissions).Error
}

// ReplacePermissions replace direct permissions of user.
// @param uint
// @param collections.Permission
// @return error
func (repository *UserRepository) ReplacePermissions(userID uint, permissions collections.Permission) error {
	return repository.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_permissions.user_id = ?", userID).Delete(&pivot.UserPermissions{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		var userPermissions []pivot.UserPermissions
		for _, permission := range permissions.Origin() {
			userPermissions = append(userPermissions, pivot.UserPermissions{
				UserID:       userID,
				PermissionID: permission.ID,
			})
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&userPermissions).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
}

// RemovePermissions remove direct permissions of user.
// @param uint
// @param collections.Permission
// @return error
func (repository *UserRepository) RemovePermissions(userID uint, permissions collections.Permission) error {
	var userPermissions []pivot.UserPermissions
	for _, permission := range permissions.Origin() {
		userPermissions = append(userPermissions, pivot.UserPermissions{
			UserID:       userID,
			PermissionID: permission.ID,
		})
	}
	return repository.Database.Delete(&userPermissions).Error
}

// ClearPermissions remove all direct permissions of user.
// @param uint
// @return error
func (repository *UserRepository) ClearPermissions(userID uint) (err error) {
	return repository.Database.Where("user_permissions.user_id = ?", userID).Delete(&pivot.UserPermissions{}).Error
}

// AddRoles add roles to user.
// @param uint
// @param collections.Role
// @return error
func (repository *UserRepository) AddRoles(userID uint, roles collections.Role) error {
	var userRoles []pivot.UserRoles
	for _, role := range roles.Origin() {
		userRoles = append(userRoles, pivot.UserRoles{
			UserID: userID,
			RoleID: role.ID,
		})
	}
	return repository.Database.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRoles).Error
}

// ReplaceRoles replace roles of user.
// @param uint
// @param collections.Role
// @return error
func (repository *UserRepository) ReplaceRoles(userID uint, roles collections.Role) error {
	return repository.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_roles.user_id = ?", userID).Delete(&pivot.UserRoles{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		var userRoles []pivot.UserRoles
		for _, role := range roles.Origin() {
			userRoles = append(userRoles, pivot.UserRoles{
				UserID: userID,
				RoleID: role.ID,
			})
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&userRoles).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

// RemoveRoles remove roles of user.
// @param uint
// @param collections.Role
// @return error
func (repository *UserRepository) RemoveRoles(userID uint, roles collections.Role) error {
	var userRoles []pivot.UserRoles
	for _, role := range roles.Origin() {
		userRoles = append(userRoles, pivot.UserRoles{
			UserID: userID,
			RoleID: role.ID,
		})
	}
	return repository.Database.Delete(&userRoles).Error
}

// ClearRoles remove all roles of user.
// @param uint
// @return error
func (repository *UserRepository) ClearRoles(userID uint) (err error) {
	return repository.Database.Where("user_roles.user_id = ?", userID).Delete(&pivot.UserRoles{}).Error
}

// CONTROLS

// HasRole does the user have the given role?
// @param uint
// @param models.Role
// @return bool, error
func (repository *UserRepository) HasRole(userID uint, role models.Role) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Where("user_roles.role_id = ?", role.ID).Count(&count).Error
	return count > 0, err
}

// HasAllRoles does the user have all the given roles?
// @param uint
// @param collections.Role
// @return bool, error
func (repository *UserRepository) HasAllRoles(userID uint, roles collections.Role) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Where("user_roles.role_id IN (?)", roles.IDs()).Count(&count).Error
	return roles.Len() == count, err
}

// HasAnyRoles does the user have any of the given roles?
// @param uint
// @param collections.Role
// @return bool, error
func (repository *UserRepository) HasAnyRoles(userID uint, roles collections.Role) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Where("user_roles.role_id IN (?)", roles.IDs()).Count(&count).Error
	return count > 0, err
}

// HasDirectPermission does the user have the given permission? (not including the permissions of the roles)
// @param uint
// @param collections.Permission
// @return bool, error
func (repository *UserRepository) HasDirectPermission(userID uint, permission models.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Where("user_permissions.permission_id = ?", permission.ID).Count(&count).Error
	return count > 0, err
}

// HasAllDirectPermissions does the user have all the given permissions? (not including the permissions of the roles)
// @param uint
// @param collections.Permission
// @return bool, error
func (repository *UserRepository) HasAllDirectPermissions(userID uint, permissions collections.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Where("user_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return permissions.Len() == count, err
}

// HasAnyDirectPermissions does the user have any of the given permissions? (not including the permissions of the roles)
// @param uint
// @param collections.Permission
// @return bool, error
func (repository *UserRepository) HasAnyDirectPermissions(userID uint, permissions collections.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Where("user_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return count > 0, err
}
