package repositories

import (
	"gorm.io/gorm"

	"github.com/Permify/permify-gorm/collections"
	"github.com/Permify/permify-gorm/models"
	"github.com/Permify/permify-gorm/models/pivot"
	"github.com/Permify/permify-gorm/repositories/scopes"
)

// IPermissionRepository its data access layer abstraction of permission.
type IPermissionRepository interface {
	Migratable

	// single fetch options

	GetPermissionByID(ID uint) (permission models.Permission, err error)
	GetPermissionByGuardName(guardName string) (permission models.Permission, err error)

	// Multiple fetch options

	GetPermissions(IDs []uint) (permissions collections.Permission, err error)
	GetPermissionsByGuardNames(guardNames []string) (permissions collections.Permission, err error)

	// ID fetch options

	GetPermissionIDs(pagination scopes.GormPager) (permissionIDs []uint, totalCount int64, err error)
	GetDirectPermissionIDsOfUserByID(userID uint, pagination scopes.GormPager) (permissionIDs []uint, totalCount int64, err error)
	GetPermissionIDsOfRolesByIDs(roleIDs []uint, pagination scopes.GormPager) (permissionIDs []uint, totalCount int64, err error)

	// FirstOrCreate & Updates & Delete

	FirstOrCreate(permission *models.Permission) (err error)
	Updates(permission *models.Permission, updates map[string]interface{}) (err error)
	Delete(permission *models.Permission) (err error)
}

// PermissionRepository its data access layer of permission.
type PermissionRepository struct {
	Database *gorm.DB
}

// Migrate generate tables from the database.
// @return error
func (repository *PermissionRepository) Migrate() (err error) {
	err = repository.Database.AutoMigrate(models.Permission{})
	err = repository.Database.AutoMigrate(pivot.UserPermissions{})
	return
}

// GetPermissionByID get permission by id.
// @param uint
// @return models.Permission, error
func (repository *PermissionRepository) GetPermissionByID(ID uint) (permission models.Permission, err error) {
	err = repository.Database.First(&permission, "permissions.id = ?", ID).Error
	return
}

// GetPermissionByGuardName get permission by guard name.
// @param string
// @return models.Permission, error
func (repository *PermissionRepository) GetPermissionByGuardName(guardName string) (permission models.Permission, err error) {
	err = repository.Database.Where("permissions.guard_name = ?", guardName).First(&permission).Error
	return
}

// MULTIPLE FETCH OPTIONS

// GetPermissions get permissions by ids.
// @param []uint
// @return collections.Role, error
func (repository *PermissionRepository) GetPermissions(IDs []uint) (permissions collections.Permission, err error) {
	err = repository.Database.Where("permissions.id IN (?)", IDs).Find(&permissions).Error
	return
}

// GetPermissionsByGuardNames get permissions by guard names.
// @param []string
// @return collections.Permission, error
func (repository *PermissionRepository) GetPermissionsByGuardNames(guardNames []string) (permissions collections.Permission, err error) {
	err = repository.Database.Where("permissions.guard_name IN (?)", guardNames).Find(&permissions).Error
	return
}

// ID FETCH OPTIONS

// GetPermissionIDs get permission ids. (with pagination)
// @param repositories_scopes.GormPager
// @return []uint, int64, error
func (repository *PermissionRepository) GetPermissionIDs(pagination scopes.GormPager) (permissionIDs []uint, totalCount int64, err error) {
	err = repository.Database.Model(&models.Permission{}).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("permissions.id", &permissionIDs).Error
	return
}

// GetDirectPermissionIDsOfUserByID get direct permission ids of user. (with pagination)
// @param uint
// @param repositories_scopes.GormPager
// @return []uint, int64, error
func (repository *PermissionRepository) GetDirectPermissionIDsOfUserByID(userID uint, pagination scopes.GormPager) (permissionIDs []uint, totalCount int64, err error) {
	err = repository.Database.Table("user_permissions").Where("user_permissions.user_id = ?", userID).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("user_permissions.permission_id", &permissionIDs).Error
	return
}

// GetPermissionIDsOfRolesByIDs get permission ids of roles. (with pagination)
// @param []uint
// @param repositories_scopes.GormPager
// @return []uint, int64, error
func (repository *PermissionRepository) GetPermissionIDsOfRolesByIDs(roleIDs []uint, pagination scopes.GormPager) (permissionIDs []uint, totalCount int64, err error) {
	err = repository.Database.Table("role_permissions").Distinct("role_permissions.permission_id").Where("role_permissions.role_id IN (?)", roleIDs).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("role_permissions.permission_id", &permissionIDs).Error
	return
}

// FirstOrCreate & Updates & Delete

// FirstOrCreate create new permission if name not exist.
// @param *models.Permission
// @return error
func (repository *PermissionRepository) FirstOrCreate(permission *models.Permission) error {
	return repository.Database.Where(models.Role{GuardName: permission.GuardName}).FirstOrCreate(permission).Error
}

// Updates update permission.
// @param *models.Permission
// @param map[string]interface{}
// @return error
func (repository *PermissionRepository) Updates(permission *models.Permission, updates map[string]interface{}) (err error) {
	return repository.Database.Model(permission).Updates(updates).Error
}

// Delete delete permission.
// @param *models.Permission
// @return error
func (repository *PermissionRepository) Delete(permission *models.Permission) (err error) {
	return repository.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_permissions.permission_id = ?", permission.ID).Delete(&pivot.UserPermissions{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Delete(permission).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
}

// paginate pagging if pagination option is true.
// @param repositories_scopes.GormPager
// @return func(db *gorm.DB) *gorm.DB
func (repository *PermissionRepository) paginate(pagination scopes.GormPager) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination != nil {
			db.Scopes(pagination.ToPaginate())
		}

		return db
	}
}
