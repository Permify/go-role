package repositories

import (
	`gorm.io/gorm`

	`github.com/Permify/permify-gorm/collections`
	`github.com/Permify/permify-gorm/models`
	models_pivot `github.com/Permify/permify-gorm/models/models.pivot`
	repositories_scopes `github.com/Permify/permify-gorm/repositories/repositories.scopes`
)

// IRoleRepository
// Its data access layer abstraction of role
type IRoleRepository interface {
	Migratable

	// single fetch options

	GetRoleByID(ID uint) (role models.Role, err error)
	GetRoleByIDWithPermissions(ID uint) (role models.Role, err error)

	GetRoleByGuardName(guardName string) (role models.Role, err error)
	GetRoleByGuardNameWithPermissions(guardName string) (role models.Role, err error)

	// Multiple fetch options

	GetRoles(roleIDs []uint) (roles collections.Role, err error)
	GetRolesWithPermissions(roleIDs []uint) (roles collections.Role, err error)

	GetRolesByGuardNames(guardNames []string) (roles collections.Role, err error)
	GetRolesByGuardNamesWithPermissions(guardNames []string) (roles collections.Role, err error)

	// ID fetch options

	GetRoleIDs(pagination repositories_scopes.GormPager) (roleIDs []uint, totalCount int64, err error)
	GetRoleIDsOfUser(userID uint, pagination repositories_scopes.GormPager) (roleIDs []uint, totalCount int64, err error)
	GetRoleIDsOfPermission(permissionID uint, pagination repositories_scopes.GormPager) (roleIDs []uint, totalCount int64, err error)

	// FirstOrCreate & Updates & Delete

	FirstOrCreate(role *models.Role) (err error)
	Updates(role *models.Role, updates map[string]interface{}) (err error)
	Delete(role *models.Role) (err error)

	// Actions

	AddPermissions(role *models.Role, permissions collections.Permission) (err error)
	ReplacePermissions(role *models.Role, permissions collections.Permission) (err error)
	RemovePermissions(role *models.Role, permissions collections.Permission) (err error)
	ClearPermissions(role *models.Role) (err error)

	// Controls

	HasPermission(roles collections.Role, permission models.Permission) (b bool, err error)
	HasAllPermissions(roles collections.Role, permissions collections.Permission) (b bool, err error)
	HasAnyPermissions(roles collections.Role, permissions collections.Permission) (b bool, err error)
}

// RoleRepository
// Its data access layer of role
type RoleRepository struct {
	Database *gorm.DB
}

// Migrate
// Generate tables from the database
// @return error
func (repository *RoleRepository) Migrate() (err error) {
	err = repository.Database.AutoMigrate(models.Role{})
	err = repository.Database.AutoMigrate(models_pivot.UserRoles{})
	return
}

// SINGLE FETCH OPTIONS

// GetRoleByID
// Get role by id.
// @param uint
// @return models.Role, error
func (repository *RoleRepository) GetRoleByID(ID uint) (role models.Role, err error) {
	err = repository.Database.First(&role, "roles.id = ?", ID).Error
	return
}

// GetRoleByIDWithPermissions
// Get role by id with its permissions.
// @param uint
// @return models.Role, error
func (repository *RoleRepository) GetRoleByIDWithPermissions(ID uint) (role models.Role, err error) {
	err = repository.Database.Preload("Permissions").First(&role, "roles.id = ?", ID).Error
	return
}

// GetRoleByGuardName
// Get role by guard name.
// @param string
// @return models.Role, error
func (repository *RoleRepository) GetRoleByGuardName(guardName string) (role models.Role, err error) {
	err = repository.Database.Where("roles.guard_name = ?", guardName).First(&role).Error
	return
}

// GetRoleByGuardNameWithPermissions
// Get role by guard name with its permissions.
// @param string
// @return models.Role, error
func (repository *RoleRepository) GetRoleByGuardNameWithPermissions(guardName string) (role models.Role, err error) {
	err = repository.Database.Preload("Permissions").Where("roles.guard_name = ?", guardName).First(&role).Error
	return
}

// MULTIPLE FETCH OPTIONS

// GetRoles
// Get roles by ids.
// @param []uint
// @return collections.Role, error
func (repository *RoleRepository) GetRoles(IDs []uint) (roles collections.Role, err error) {
	err = repository.Database.Where("roles.id IN (?)", IDs).Find(&roles).Error
	return
}

// GetRolesWithPermissions
// Get roles by ids with its permissions.
// @param []uint
// @return collections.Role, error
func (repository *RoleRepository) GetRolesWithPermissions(IDs []uint) (roles collections.Role, err error) {
	err = repository.Database.Preload("Permissions").Where("roles.id IN (?)", IDs).Find(&roles).Error
	return
}

// GetRolesByGuardNames
// Get roles by guard names.
// @param []string
// @return collections.Role, error
func (repository *RoleRepository) GetRolesByGuardNames(guardNames []string) (roles collections.Role, err error) {
	err = repository.Database.Where("roles.guard_name IN (?)", guardNames).Find(&roles).Error
	return
}

// GetRolesByGuardNamesWithPermissions
// Get roles by guard names.
// @param []string
// @return collections.Role, error
func (repository *RoleRepository) GetRolesByGuardNamesWithPermissions(guardNames []string) (roles collections.Role, err error) {
	err = repository.Database.Preload("Permissions").Where("roles.guard_name IN (?)", guardNames).Find(&roles).Error
	return
}

// ID FETCH OPTIONS

// GetRoleIDs
// Get role ids. (with pagination)
// @param repositories_scopes.GormPager
// @return []uint, int64, error
func (repository *RoleRepository) GetRoleIDs(pagination repositories_scopes.GormPager) (roleIDs []uint, totalCount int64, err error) {
	err = repository.Database.Model(&models.Role{}).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("roles.id", &roleIDs).Error
	return
}

// GetRoleIDsOfUser
// Get role ids of user. (with pagination)
// @param uint
// @param repositories_scopes.GormPager
// @return []uint, int64, error
func (repository *RoleRepository) GetRoleIDsOfUser(userID uint, pagination repositories_scopes.GormPager) (roleIDs []uint, totalCount int64, err error) {
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("user_roles.role_id", &roleIDs).Error
	return
}

// GetRoleIDsOfPermission
// Get role ids of permission. (with pagination)
// @param uint
// @param repositories_scopes.GormPager
// @return []uint, int64, error
func (repository *RoleRepository) GetRoleIDsOfPermission(permissionID uint, pagination repositories_scopes.GormPager) (roleIDs []uint, totalCount int64, err error) {
	err = repository.Database.Table("role_permissions").Where("role_permissions.permission_id = ?", permissionID).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("role_permissions.role_id", &roleIDs).Error
	return
}

// FirstOrCreate & Updates & Delete

// FirstOrCreate
// Create new role if name not exist.
// @param *models.Role
// @return error
func (repository *RoleRepository) FirstOrCreate(role *models.Role) error {
	return repository.Database.Where(models.Role{GuardName: role.GuardName}).FirstOrCreate(role).Error
}

// Updates
// Update role.
// @param *models.Role
// @param map[string]interface{}
// @return error
func (repository *RoleRepository) Updates(role *models.Role, updates map[string]interface{}) (err error) {
	return repository.Database.Model(role).Updates(updates).Error
}

// Delete
// Delete role.
// @param *models.Role
// @return error
func (repository *RoleRepository) Delete(role *models.Role) (err error) {
	return repository.Database.Delete(role).Error
}

// ACTIONS

// AddPermissions
// Add permissions to role.
// @param *models.Role
// @param collections.Permission
// @return error
func (repository *RoleRepository) AddPermissions(role *models.Role, permissions collections.Permission) error {
	return repository.Database.Model(role).Association("Permissions").Append(permissions.Origin())
}

// ReplacePermissions
// Replace permissions of role.
// @param *models.Role
// @param collections.Permission
// @return error
func (repository *RoleRepository) ReplacePermissions(role *models.Role, permissions collections.Permission) error {
	return repository.Database.Model(role).Association("Permissions").Replace(permissions.Origin())
}

// RemovePermissions
// Remove permissions from role.
// @param *models.Role
// @param collections.Permission
// @return error
func (repository *RoleRepository) RemovePermissions(role *models.Role, permissions collections.Permission) error {
	return repository.Database.Model(role).Association("Permissions").Delete(permissions.Origin())
}

// ClearPermissions
// Remove all permissions from role.
// @param *models.Role
// @return error
func (repository *RoleRepository) ClearPermissions(role *models.Role) (err error) {
	return repository.Database.Model(role).Association("Permissions").Clear()
}

// Controls

// HasPermission
// Does the role or any of the roles have given permission?
// @param collections.Role
// @param models.Permission
// @return bool, error
func (repository *RoleRepository) HasPermission(roles collections.Role, permission models.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("role_permissions").Where("role_permissions.role_id IN (?)", roles.IDs()).Where("role_permissions.permission_id = ?", permission).Count(&count).Error
	return count > 0, err
}

// HasAllPermissions
// Does the role or roles have all the given permissions?
// @param collections.Role
// @param collections.Permission
// @return bool, error
func (repository *RoleRepository) HasAllPermissions(roles collections.Role, permissions collections.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("role_permissions").Where("role_permissions.role_id IN (?)", roles.IDs()).Where("role_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return roles.Len()*permissions.Len() == count, err
}

// HasAnyPermissions
// Does the role or roles have any of the given permissions?
// @param collections.Role
// @param collections.Permission
// @return bool, error
func (repository *RoleRepository) HasAnyPermissions(roles collections.Role, permissions collections.Permission) (b bool, err error) {
	var count int64
	err = repository.Database.Table("role_permissions").Where("role_permissions.role_id IN (?)", roles.IDs()).Where("role_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return count > 0, err
}

/**
 * Scopes
 *
 */

// paginate
// @param repositories_scopes.GormPager
// @return func(db *gorm.DB) *gorm.DB
func (repository *RoleRepository) paginate(pagination repositories_scopes.GormPager) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if pagination != nil {
			db.Scopes(pagination.ToPaginate())
		}

		return db
	}
}
