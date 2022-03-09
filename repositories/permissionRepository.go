package repositories

import (
	`gorm.io/gorm`

	`github.com/Permify/permify-gorm/collections`
	`github.com/Permify/permify-gorm/models`
	models_pivot `github.com/Permify/permify-gorm/models/models.pivot`
	repositories_scopes `github.com/Permify/permify-gorm/repositories/repositories.scopes`
)

type IPermissionRepository interface {
	Migratable

	// single fetch options
	GetPermissionByID(ID uint) (permission models.Permission, err error)
	GetPermissionByGuardName(guardName string) (permission models.Permission, err error)

	// Multiple fetch options
	GetPermissions(IDs []uint) (permissions collections.Permission, err error)
	GetPermissionsByGuardNames(guardNames []string) (permissions collections.Permission, err error)

	// ID fetch options
	GetPermissionIDs(pagination repositories_scopes.GormPager) (permissionIDs []uint, totalCount int64, err error)
	GetDirectPermissionIDsOfUserByID(userID uint, pagination repositories_scopes.GormPager) (permissionIDs []uint, totalCount int64, err error)
	GetPermissionIDsOfRolesByIDs(roleIDs []uint, pagination repositories_scopes.GormPager) (permissionIDs []uint, totalCount int64, err error)

	// FirstOrCreate & Updates & Delete
	FirstOrCreate(permission *models.Permission) (err error)
	Updates(permission *models.Permission, updates map[string]interface{}) (err error)
	Delete(permission *models.Permission) (err error)
}

type PermissionRepository struct {
	Database *gorm.DB
}

/**
 * Migrate
 *
 * @return error
 */

func (repository *PermissionRepository) Migrate() (err error) {
	err = repository.Database.AutoMigrate(models.Permission{})
	err = repository.Database.AutoMigrate(models_pivot.UserPermissions{})
	return
}

/**
 * Single fetch options
 *
 */

func (repository *PermissionRepository) GetPermissionByID(ID uint) (permission models.Permission, err error) {
	err = repository.Database.First(&permission, "permissions.id = ?", ID).Error
	return
}

func (repository *PermissionRepository) GetPermissionByGuardName(guardName string) (permission models.Permission, err error) {
	err = repository.Database.Where("permissions.guard_name = ?", guardName).First(&permission).Error
	return
}

/**
 * Multiple fetch options
 *
 */

func (repository *PermissionRepository) GetPermissions(IDs []uint) (permissions collections.Permission, err error) {
	err = repository.Database.Where("permissions.id IN (?)", IDs).Find(&permissions).Error
	return
}

func (repository *PermissionRepository) GetPermissionsByGuardNames(guardNames []string) (permissions collections.Permission, err error) {
	err = repository.Database.Where("permissions.guard_name IN (?)", guardNames).Find(&permissions).Error
	return
}

/**
 * ID fetch options
 *
 */

func (repository *PermissionRepository) GetPermissionIDs(pagination repositories_scopes.GormPager) (permissionIDs []uint, totalCount int64, err error) {
	err = repository.Database.Model(&models.Permission{}).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("permissions.id", &permissionIDs).Error
	return
}

func (repository *PermissionRepository) GetDirectPermissionIDsOfUserByID(userID uint, pagination repositories_scopes.GormPager) (permissionIDs []uint, totalCount int64, err error) {
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("user_permissions.permission_id", &permissionIDs).Error
	return
}

func (repository *PermissionRepository) GetPermissionIDsOfRolesByIDs(roleIDs []uint, pagination repositories_scopes.GormPager) (permissionIDs []uint, totalCount int64, err error) {
	err = repository.Database.Table("role_permissions").Distinct("role_permissions.permission_id").Where("role_permissions.role_id IN (?)", roleIDs).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("role_permissions.permission_id", &permissionIDs).Error
	return
}

/**
 * FirstOrCreate & Updates & Delete
 *
 */

func (repository *PermissionRepository) FirstOrCreate(permission *models.Permission) error {
	return repository.Database.Where(models.Role{GuardName: permission.GuardName}).FirstOrCreate(permission).Error
}

func (repository *PermissionRepository) Updates(permission *models.Permission, updates map[string]interface{}) (err error) {
	return repository.Database.Model(permission).Updates(updates).Error
}

func (repository *PermissionRepository) Delete(permission *models.Permission) (err error) {
	return repository.Database.Delete(permission).Error
}

/**
 * Scopes
 *
 */

func (repository *PermissionRepository) paginate(pagination repositories_scopes.GormPager) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if pagination != nil {
			db.Scopes(pagination.ToPaginate())
		}

		return db
	}
}
