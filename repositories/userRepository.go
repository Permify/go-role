package repositories

import (
	`gorm.io/gorm`
	`gorm.io/gorm/clause`

	`github.com/Permify/permify-gorm/collections`
	`github.com/Permify/permify-gorm/models`
	models_pivot `github.com/Permify/permify-gorm/models/models.pivot`
)

type IUserRepository interface {

	// permission
	AddPermissions(userID uint, permissions collections.Permission) (err error)
	ReplacePermissions(userID uint, permissions collections.Permission) (err error)
	RemovePermissions(userID uint, permissions collections.Permission) (err error)
	ClearPermissions(userID uint) (err error)

	// role
	AddRoles(userID uint, roles collections.Role) (err error)
	ReplaceRoles(userID uint, roles collections.Role) (err error)
	RemoveRoles(userID uint, roles collections.Role) (err error)
	ClearRoles(userID uint) (err error)

	// Controls
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

/**
 * Permission
 *
 */

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

func (repository *UserRepository) ClearPermissions(userID uint) (err error) {
	return repository.Database.Where("user_permissions.user_id = ?", userID).Delete(&models_pivot.UserPermissions{}).Error
}

/**
 * Role
 *
 */

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

func (repository *UserRepository) ClearRoles(userID uint) (err error) {
	return repository.Database.Where("user_roles.user_id = ?", userID).Delete(&models_pivot.UserRoles{}).Error
}

// Controls

// role

func (repository *UserRepository) HasRole(userID uint, role models.Role) (b bool, err error)  {
	var count int64
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Where("user_roles.role_id = ?", role.ID).Count(&count).Error
	return count > 0, err
}

func (repository *UserRepository) HasAllRoles(userID uint, roles collections.Role) (b bool, err error)  {
	var count int64
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Where("user_roles.role_id IN (?)", roles.IDs()).Count(&count).Error
	return roles.Len() == count, err
}

func (repository *UserRepository) HasAnyRoles(userID uint, roles collections.Role) (b bool, err error)  {
	var count int64
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Where("user_roles.role_id IN (?)", roles.IDs()).Count(&count).Error
	return count > 0, err
}

// permission

func (repository *UserRepository) HasDirectPermission(userID uint, permission models.Permission) (b bool,err error)  {
	var count int64
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Where("user_permissions.permission_id = ?", permission.ID).Count(&count).Error
	return count > 0, err
}

func (repository *UserRepository) HasAllDirectPermissions(userID uint, permissions collections.Permission) (b bool, err error)  {
	var count int64
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Where("user_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return permissions.Len() == count, err
}

func (repository *UserRepository) HasAnyDirectPermissions(userID uint, permissions collections.Permission) (b bool, err error)  {
	var count int64
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Where("user_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return count > 0, err
}
