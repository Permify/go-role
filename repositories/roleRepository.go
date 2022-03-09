package repositories

import (
	`gorm.io/gorm`

	`github.com/Permify/permify-gorm/collections`
	`github.com/Permify/permify-gorm/models`
	models_pivot `github.com/Permify/permify-gorm/models/models.pivot`
	repositories_scopes `github.com/Permify/permify-gorm/repositories/repositories.scopes`
)

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

type RoleRepository struct {
	Database *gorm.DB
}

/**
 * Migrate
 *
 * @return error
 */

func (repository *RoleRepository) Migrate() (err error) {
	err = repository.Database.AutoMigrate(models.Role{})
	err = repository.Database.AutoMigrate(models_pivot.UserRoles{})
	return
}

/**
 * Single fetch options
 *
 */

func (repository *RoleRepository) GetRoleByID(ID uint) (role models.Role, err error) {
	err = repository.Database.First(&role, "roles.id = ?", ID).Error
	return
}

func (repository *RoleRepository) GetRoleByIDWithPermissions(ID uint) (role models.Role, err error) {
	err = repository.Database.Preload("Permissions").First(&role, "roles.id = ?", ID).Error
	return
}

func (repository *RoleRepository) GetRoleByGuardName(guardName string) (role models.Role, err error) {
	err = repository.Database.Where("roles.guard_name = ?", guardName).First(&role).Error
	return
}

func (repository *RoleRepository) GetRoleByGuardNameWithPermissions(guardName string) (role models.Role, err error) {
	err = repository.Database.Preload("Permissions").Where("roles.guard_name = ?", guardName).First(&role).Error
	return
}

/**
 * Multiple fetch options
 *
 */

func (repository *RoleRepository) GetRoles(IDs []uint) (roles collections.Role, err error) {
	err = repository.Database.Where("roles.id IN (?)", IDs).Find(&roles).Error
	return
}

func (repository *RoleRepository) GetRolesWithPermissions(IDs []uint) (roles collections.Role, err error) {
	err = repository.Database.Preload("Permissions").Where("roles.id IN (?)", IDs).Find(&roles).Error
	return
}


func (repository *RoleRepository) GetRolesByGuardNames(guardNames []string) (roles collections.Role, err error) {
	err = repository.Database.Where("roles.guard_name IN (?)", guardNames).Find(&roles).Error
	return
}

func (repository *RoleRepository) GetRolesByGuardNamesWithPermissions(guardNames []string) (roles collections.Role, err error) {
	err = repository.Database.Preload("Permissions").Where("roles.guard_name IN (?)", guardNames).Find(&roles).Error
	return
}


/**
 * ID fetch options
 *
 */

func (repository *RoleRepository) GetRoleIDs(pagination repositories_scopes.GormPager) (roleIDs []uint, totalCount int64, err error) {
	err = repository.Database.Model(&models.Role{}).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("roles.id", &roleIDs).Error
	return
}

func (repository *RoleRepository) GetRoleIDsOfUser(userID uint, pagination repositories_scopes.GormPager) (roleIDs []uint, totalCount int64, err error) {
	err = repository.Database.Table("user_roles").Where("user_roles.user_id = ?", userID).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("user_roles.role_id", &roleIDs).Error
	return
}

func (repository *RoleRepository) GetRoleIDsOfPermission(permissionID uint, pagination repositories_scopes.GormPager) (roleIDs []uint, totalCount int64, err error) {
	err = repository.Database.Table("role_permissions").Where("role_permissions.permission_id = ?", permissionID).Count(&totalCount).Scopes(repository.paginate(pagination)).Pluck("role_permissions.role_id", &roleIDs).Error
	return
}

/**
 * FirstOrCreate & Updates & Delete
 *
 */

func (repository *RoleRepository) FirstOrCreate(role *models.Role) error {
	return repository.Database.Where(models.Role{GuardName: role.GuardName}).FirstOrCreate(role).Error
}

func (repository *RoleRepository) Updates(role *models.Role, updates map[string]interface{}) (err error) {
	return repository.Database.Model(role).Updates(updates).Error
}

func (repository *RoleRepository) Delete(role *models.Role) (err error) {
	return repository.Database.Delete(role).Error
}

/**
 * Actions
 *
 */

func (repository *RoleRepository) AddPermissions(role *models.Role, permissions collections.Permission) error {
	return repository.Database.Model(role).Association("Permissions").Append(permissions.Origin())
}

func (repository *RoleRepository) ReplacePermissions(role *models.Role, permissions collections.Permission) error {
	return repository.Database.Model(role).Association("Permissions").Replace(permissions.Origin())
}

func (repository *RoleRepository) RemovePermissions(role *models.Role, permissions collections.Permission) error {
	return repository.Database.Model(role).Association("Permissions").Delete(permissions.Origin())
}

func (repository *RoleRepository) ClearPermissions(role *models.Role) (err error) {
	return repository.Database.Model(role).Association("Permissions").Clear()
}

// Controls

func (repository *RoleRepository) HasPermission(roles collections.Role, permission models.Permission) (b bool, err error)  {
	var count int64
	err = repository.Database.Table("role_permissions").Where("role_permissions.role_id IN (?)", roles.IDs()).Where("role_permissions.permission_id = ?", permission).Count(&count).Error
	return count > 0, err
}

func (repository *RoleRepository) HasAllPermissions(roles collections.Role, permissions collections.Permission) (b bool, err error)  {
	var count int64
	err = repository.Database.Table("role_permissions").Where("role_permissions.role_id IN (?)", roles.IDs()).Where("role_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return roles.Len() * permissions.Len() == count, err
}

func (repository *RoleRepository) HasAnyPermissions(roles collections.Role, permissions collections.Permission) (b bool, err error)  {
	var count int64
	err = repository.Database.Table("role_permissions").Where("role_permissions.role_id IN (?)", roles.IDs()).Where("role_permissions.permission_id IN (?)", permissions.IDs()).Count(&count).Error
	return count > 0, err
}

/**
 * Scopes
 *
 */

func (repository *RoleRepository) paginate(pagination repositories_scopes.GormPager) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if pagination != nil {
			db.Scopes(pagination.ToPaginate())
		}

		return db
	}
}
