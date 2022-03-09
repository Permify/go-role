package permify_gorm

import (
	`errors`
	`gorm.io/gorm`

	`github.com/Permify/permify-gorm/collections`
	`github.com/Permify/permify-gorm/helpers`
	`github.com/Permify/permify-gorm/models`
	`github.com/Permify/permify-gorm/options`
	`github.com/Permify/permify-gorm/repositories`
	repositories_scopes `github.com/Permify/permify-gorm/repositories/repositories.scopes`
)

var (
	errUnsupportedValueType     = errors.New("err unsupported value type")
)

type Options struct {
	Migrate bool
	DB      *gorm.DB
}

// New initializer for Permify
// If migration is true, it generate all tables in the database if they don't exist.
func New(opts Options) (p *Permify, err error) {

	roleRepository := &repositories.RoleRepository{Database: opts.DB}
	permissionRepository := &repositories.PermissionRepository{Database: opts.DB}
	userRepository := &repositories.UserRepository{Database: opts.DB}

	if opts.Migrate {
		err = repositories.Migrates(roleRepository, permissionRepository)
		if err != nil {
			return nil, err
		}
	}

	p = &Permify{
		RoleRepository:       roleRepository,
		PermissionRepository: permissionRepository,
		UserRepository:       userRepository,
	}

	return
}

type Permify struct {
	RoleRepository       repositories.IRoleRepository
	PermissionRepository repositories.IPermissionRepository
	UserRepository       repositories.IUserRepository
}

// ROLE

// GetRole
// Fetch role according to the role name or id.
// If withPermissions is true, it will preload the permissions to the role.
// If the given variable is an array, the first element of the given array is returned.
func (s *Permify) GetRole(r interface{}, withPermissions bool) (role models.Role, err error) {

	if helpers.IsArray(r) {
		var roles []models.Role
		roles, err = s.GetRoles(r, withPermissions)
		if err != nil {
			return models.Role{}, err
		}
		if len(roles) > 0 {
			role = roles[0]
		}
		return
	}

	if helpers.IsString(r) {
		if withPermissions {
			return s.RoleRepository.GetRoleByGuardNameWithPermissions(helpers.Guard(r.(string)))
		}
		return s.RoleRepository.GetRoleByGuardName(helpers.Guard(r.(string)))
	}

	if helpers.IsInt(r) {
		if withPermissions {
			return s.RoleRepository.GetRoleByIDWithPermissions(uint(r.(int)))
		}
		return s.RoleRepository.GetRoleByID(uint(r.(int)))
	}

	return models.Role{}, errUnsupportedValueType
}

// GetRoles
// Fetch roles according to the role names or ids.
// If withPermissions is true, it will preload the permissions to the roles.
func (s *Permify) GetRoles(r interface{}, withPermissions bool) (roles collections.Role, err error) {

	if !helpers.IsArray(r) {
		var role models.Role
		role, err = s.GetRole(r, withPermissions)
		if err != nil {
			return collections.Role{}, err
		}
		roles = collections.Role{role}
		return
	}

	if helpers.IsStringArray(r) {
		if withPermissions {
			return s.RoleRepository.GetRolesByGuardNamesWithPermissions(helpers.GuardArray(r.([]string)))
		}
		return s.RoleRepository.GetRolesByGuardNames(helpers.GuardArray(r.([]string)))
	}

	if helpers.IsUIntArray(r) {
		if withPermissions {
			return s.RoleRepository.GetRolesWithPermissions(r.([]uint))
		}
		return s.RoleRepository.GetRoles(r.([]uint))
	}

	return collections.Role{}, errUnsupportedValueType
}

// GetAllRoles */
// Fetch all the roles. (with pagination option).
// If withPermissions is true, it will preload the permissions to the role.
func (s *Permify) GetAllRoles(option options.RoleOption) (roles collections.Role, totalCount int64, err error) {
	var roleIDs []uint
	if option.Pagination == nil {
		roleIDs, totalCount, err = s.RoleRepository.GetRoleIDs(nil)
	} else {
		roleIDs, totalCount, err = s.RoleRepository.GetRoleIDs(&repositories_scopes.GormPagination{Pagination: option.Pagination.Get()})
	}
	roles, err = s.GetRoles(roleIDs, option.WithPermissions)
	return
}

// GetRolesOfUser */
// Fetch all the roles of the user. (with pagination option).
// If withPermissions is true, it will preload the permissions to the role.
func (s *Permify) GetRolesOfUser(userID uint, option options.RoleOption) (roles collections.Role, totalCount int64, err error) {
	var roleIDs []uint
	if option.Pagination == nil {
		roleIDs, totalCount, err = s.RoleRepository.GetRoleIDsOfUser(userID, nil)
	} else {
		roleIDs, totalCount, err = s.RoleRepository.GetRoleIDsOfUser(userID, &repositories_scopes.GormPagination{Pagination: option.Pagination.Get()})
	}
	roles, err = s.GetRoles(roleIDs, option.WithPermissions)
	return
}

// CreateRole */
// Create new role
// Name parameter is converted to guard name. example: senior $#% associate -> senior-associate.
// If a role with the same name has been created before, it will not create it again. (FirstOrCreate)
func (s *Permify) CreateRole(name string, description string) (err error) {
	return s.RoleRepository.FirstOrCreate(&models.Role{
		Name:        name,
		GuardName:   helpers.Guard(name),
		Description: description,
	})
}

// DeleteRole */
// Delete role
// If the role is in use, its relations from the pivot tables are deleted.
func (s *Permify) DeleteRole(r interface{}) (err error) {
	var role models.Role
	role, err = s.GetRole(r, false)
	if err != nil {
		return err
	}
	return s.RoleRepository.Delete(&role)
}

// AddPermissionsToRole */
// add permission to role.
func (s *Permify) AddPermissionsToRole(r interface{}, p interface{}) (err error) {

	var role models.Role
	role, err = s.GetRole(r, false)
	if err != nil {
		return err
	}

	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return err
	}

	if permissions.Len() > 0 {
		err = s.RoleRepository.AddPermissions(&role, permissions)
	}

	return
}

// ReplacePermissionsToRole */
// Overwrites the permissions of the role according to the permission names or ids.
func (s *Permify) ReplacePermissionsToRole(r interface{}, p interface{}) (err error) {

	var role models.Role
	role, err = s.GetRole(r, false)
	if err != nil {
		return err
	}

	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return err
	}

	if permissions.Len() > 0 {
		err = s.RoleRepository.ReplacePermissions(&role, permissions)
	}

	return
}

// RemovePermissionsFromRole */
// remove permissions from role according to the permission names or ids.
func (s *Permify) RemovePermissionsFromRole(r interface{}, p interface{}) (err error) {

	var role models.Role
	role, err = s.GetRole(r, false)
	if err != nil {
		return err
	}

	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return err
	}

	if permissions.Len() > 0 {
		err = s.RoleRepository.RemovePermissions(&role, permissions)
	}

	return
}

// PERMISSION

// GetPermission */
// Fetch permission according to the permission name or id.
// If the given variable is an array, the first element of the given array is returned.
func (s *Permify) GetPermission(p interface{}) (permission models.Permission, err error) {

	if helpers.IsArray(p) {
		var permissions []models.Permission
		permissions, err = s.GetPermissions(p)
		if err != nil {
			return models.Permission{}, err
		}
		if len(permissions) > 0 {
			permission = permissions[0]
		}
		return
	}

	if helpers.IsString(p) {
		return s.PermissionRepository.GetPermissionByGuardName(helpers.Guard(p.(string)))
	}

	if helpers.IsInt(p) {
		return s.PermissionRepository.GetPermissionByID(uint(p.(int)))
	}

	return models.Permission{}, errUnsupportedValueType
}

// GetPermissions */
// Fetch permissions according to the permission names or ids.
func (s *Permify) GetPermissions(p interface{}) (permissions collections.Permission, err error) {

	if !helpers.IsArray(p) {
		var permission models.Permission
		permission, err = s.GetPermission(p)
		if err != nil {
			return collections.Permission{}, err
		}
		permissions = collections.Permission{permission}
		return
	}

	if helpers.IsStringArray(p) {
		return s.PermissionRepository.GetPermissionsByGuardNames(helpers.GuardArray(p.([]string)))
	}

	if helpers.IsUIntArray(p) {
		return s.PermissionRepository.GetPermissions(p.([]uint))
	}

	return collections.Permission{}, errUnsupportedValueType
}

// GetAllPermissions */
// Fetch all the permissions. (with pagination option).
func (s *Permify) GetAllPermissions(option options.PermissionOption) (permissions collections.Permission, totalCount int64, err error) {
	var permissionIDs []uint
	if option.Pagination == nil {
		permissionIDs, totalCount, err = s.PermissionRepository.GetPermissionIDs(nil)
	} else {
		permissionIDs, totalCount, err = s.PermissionRepository.GetPermissionIDs(&repositories_scopes.GormPagination{Pagination: option.Pagination.Get()})
	}
	permissions, err = s.GetPermissions(permissionIDs)
	return
}

// GetDirectPermissionsOfUser */
// Fetch all direct permissions of the user. (with pagination option)
func (s *Permify) GetDirectPermissionsOfUser(userID uint, option options.PermissionOption) (permissions collections.Permission, totalCount int64, err error) {
	var permissionIDs []uint
	if option.Pagination == nil {
		permissionIDs, totalCount, err = s.PermissionRepository.GetDirectPermissionIDsOfUserByID(userID, nil)
	} else {
		permissionIDs, totalCount, err = s.PermissionRepository.GetDirectPermissionIDsOfUserByID(userID, &repositories_scopes.GormPagination{Pagination: option.Pagination.Get()})
	}
	permissions, err = s.GetPermissions(permissionIDs)
	return
}

// GetPermissionsOfRoles */
// Fetch all permissions of the roles. (with pagination option)
func (s *Permify) GetPermissionsOfRoles(r interface{}, option options.PermissionOption) (permissions collections.Permission, totalCount int64, err error) {

	var roles collections.Role
	roles, err = s.GetRoles(r, false)
	if err != nil {
		return collections.Permission{}, 0, err
	}

	var permissionIDs []uint
	if option.Pagination == nil {
		permissionIDs, totalCount, err = s.PermissionRepository.GetPermissionIDsOfRolesByIDs(roles.IDs(), nil)
	} else {
		permissionIDs, totalCount, err = s.PermissionRepository.GetPermissionIDsOfRolesByIDs(roles.IDs(), &repositories_scopes.GormPagination{Pagination: option.Pagination.Get()})
	}

	permissions, err = s.GetPermissions(permissionIDs)
	return
}

// GetAllPermissionsOfUser */
// Fetch all permissions of the user that come with direct and roles.
func (s *Permify) GetAllPermissionsOfUser(userID uint) (permissions collections.Permission, err error) {

	var urIDs []uint
	urIDs, _, err = s.RoleRepository.GetRoleIDsOfUser(userID, nil)
	if err != nil {
		return collections.Permission{}, err
	}

	var rpIDs []uint
	rpIDs, _, err = s.PermissionRepository.GetPermissionIDsOfRolesByIDs(urIDs, nil)
	if err != nil {
		return collections.Permission{}, err
	}

	var udpIDs []uint
	udpIDs, _, err = s.PermissionRepository.GetDirectPermissionIDsOfUserByID(userID, nil)
	if err != nil {
		return collections.Permission{}, err
	}

	return s.GetPermissions(helpers.RemoveDuplicateValues(helpers.JoinUintArrays(rpIDs, udpIDs)))
}

// CreatePermission */
// Create new permission
// Name parameter is converted to guard name. example: create $#% contact -> create-contact.
// If a permission with the same name has been created before, it will not create it again. (FirstOrCreate)
func (s *Permify) CreatePermission(name string, description string) (err error) {
	return s.PermissionRepository.FirstOrCreate(&models.Permission{
		Name:        name,
		GuardName:   helpers.Guard(name),
		Description: description,
	})
}

// DeletePermission */
// Delete permission
// If the permission is in use, its relations from the pivot tables are deleted.
func (s *Permify) DeletePermission(p interface{}) (err error) {
	var permission models.Permission
	permission, err = s.GetPermission(p)
	if err != nil {
		return err
	}
	return s.PermissionRepository.Delete(&permission)
}

// USER

// AddPermissionsToUser */
// Add direct permission or permissions to user according to the permission names or ids.
func (s *Permify) AddPermissionsToUser(userID uint, p interface{}) (err error) {

	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return err
	}

	if permissions.Len() > 0 {
		err = s.UserRepository.AddPermissions(userID, permissions)
	}

	return
}

// ReplacePermissionsToUser */
// Overwrites the direct permissions of the user according to the permission names or ids.
func (s *Permify) ReplacePermissionsToUser(userID uint, p interface{}) (err error) {

	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return err
	}

	if permissions.Len() > 0 {
		return s.UserRepository.ReplacePermissions(userID, permissions)
	}

	return s.UserRepository.ClearPermissions(userID)
}

// RemovePermissionsFromUser */
// remove direct permissions from user according to the permission names or ids.
func (s *Permify) RemovePermissionsFromUser(userID uint, p interface{}) (err error) {

	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return err
	}

	if permissions.Len() > 0 {
		err = s.UserRepository.RemovePermissions(userID, permissions)
	}

	return
}


// AddRolesToUser */
// Add role or roles to user according to the role names or ids.
func (s *Permify) AddRolesToUser(userID uint, r interface{}) (err error) {

	var roles collections.Role
	roles, err = s.GetRoles(r, false)
	if err != nil {
		return err
	}

	if roles.Len() > 0 {
		err = s.UserRepository.AddRoles(userID, roles)
	}

	return
}

// ReplaceRolesToUser */
// Overwrites the roles of the user according to the role names or ids.
func (s *Permify) ReplaceRolesToUser(userID uint, r interface{}) (err error) {

	var roles collections.Role
	roles, err = s.GetRoles(r, false)
	if err != nil {
		return err
	}

	if roles.Len() > 0 {
		return s.UserRepository.ReplaceRoles(userID, roles)
	}

	return s.UserRepository.ClearRoles(userID)
}

// RemoveRolesFromUser */
// remove direct permissions from user according to the role names or ids.
func (s *Permify) RemoveRolesFromUser(userID uint, r interface{}) (err error) {

	var roles collections.Role
	roles, err = s.GetRoles(r, false)
	if err != nil {
		return err
	}

	if roles.Len() > 0 {
		err = s.UserRepository.RemoveRoles(userID, roles)
	}

	return
}

// CONTROLS

// ROLE

// RoleHasPermission */
// Does the role or any of the roles have given permission?
func (s *Permify) RoleHasPermission(r interface{}, p interface{}) (b bool, err error) {

	var roles collections.Role
	roles, err = s.GetRoles(r, false)
	if err != nil {
		return false, err
	}

	var permission models.Permission
	permission, err = s.GetPermission(p)
	if err != nil {
		return false, err
	}

	return s.RoleRepository.HasPermission(roles, permission)
}

// RoleHasAllPermissions */
// Does the role or roles have all the given permissions?
func (s *Permify) RoleHasAllPermissions(r interface{}, p interface{}) (b bool, err error) {

	var roles collections.Role
	roles, err = s.GetRoles(r, false)
	if err != nil {
		return false, err
	}

	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return false, err
	}

	return s.RoleRepository.HasAllPermissions(roles, permissions)
}

// RoleHasAnyPermissions */
// Does the role or roles have any of the given permissions?
func (s *Permify) RoleHasAnyPermissions(r interface{}, p interface{}) (b bool, err error) {

	var roles collections.Role
	roles, err = s.GetRoles(r, false)
	if err != nil {
		return false, err
	}

	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return false, err
	}

	return s.RoleRepository.HasAnyPermissions(roles, permissions)
}

// USER

// UserHasRole */
// Does the user have the given role?
func (s *Permify) UserHasRole(userID uint, r interface{}) (b bool, err error) {
	var role models.Role
	role, err = s.GetRole(r, false)
	if err != nil {
		return false, err
	}
	return s.UserRepository.HasRole(userID, role)
}

// UserHasAllRoles */
// Does the user have all the given roles?
func (s *Permify) UserHasAllRoles(userID uint, r interface{}) (b bool, err error) {
	var roles collections.Role
	roles, err = s.GetRoles(r, false)
	if err != nil {
		return false, err
	}
	return s.UserRepository.HasAllRoles(userID, roles)
}

// UserHasAnyRoles */
// Does the user have any of the given roles?
func (s *Permify) UserHasAnyRoles(userID uint, r interface{}) (b bool, err error) {
	var roles collections.Role
	roles, err = s.GetRoles(r, false)
	if err != nil {
		return false, err
	}
	return s.UserRepository.HasAnyRoles(userID, roles)
}

// UserHasDirectPermission */
// Does the user have the given permission? (not including the permissions of the roles)
func (s *Permify) UserHasDirectPermission(userID uint, p interface{}) (b bool, err error) {
	var permission models.Permission
	permission, err = s.GetPermission(p)
	if err != nil {
		return false, err
	}
	return s.UserRepository.HasDirectPermission(userID, permission)
}

// UserHasAllDirectPermissions */
// Does the user have all the given permissions? (not including the permissions of the roles)
func (s *Permify) UserHasAllDirectPermissions(userID uint, p interface{}) (b bool, err error) {
	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return false, err
	}
	return s.UserRepository.HasAllDirectPermissions(userID, permissions)
}

// UserHasAnyDirectPermissions */
// Does the user have any of the given permissions? (not including the permissions of the roles)
func (s *Permify) UserHasAnyDirectPermissions(userID uint, p interface{}) (b bool, err error) {
	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return false, err
	}
	return s.UserRepository.HasAnyDirectPermissions(userID, permissions)
}


// UserHasPermission */
// Does the user have the given permission? (including the permissions of the roles)
func (s *Permify) UserHasPermission(userID uint, p interface{}) (b bool, err error) {
	var permission models.Permission
	permission, err = s.GetPermission(p)
	if err != nil {
		return false, err
	}

	var hasDirect bool
	hasDirect, err = s.UserRepository.HasDirectPermission(userID, permission)
	if err != nil {
		return false, err
	}
	if hasDirect {
		return true, err
	}

	var roleIDs []uint
	roleIDs, _, err = s.RoleRepository.GetRoleIDsOfUser(userID, nil)
	if err != nil {
		return false, err
	}

	var hasAny bool
	hasAny, err = s.RoleHasAnyPermissions(roleIDs, p)
	if err != nil {
		return false, err
	}
	if hasAny {
		return true, err
	}

	return false, err
}

// UserHasAllPermissions */
// Does the user have all the given permissions? (including the permissions of the roles)
func (s *Permify) UserHasAllPermissions(userID uint, p interface{}) (b bool, err error) {

	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return false, err
	}

	var upIDs []uint
	upIDs,_, err = s.PermissionRepository.GetDirectPermissionIDsOfUserByID(userID, nil)
	if err != nil {
		return false, err
	}

	var rIDs []uint
	rIDs,_, err = s.RoleRepository.GetRoleIDsOfUser(userID, nil)
	if err != nil {
		return false, err
	}

	var pIDs []uint
	pIDs, _, err = s.PermissionRepository.GetPermissionIDsOfRolesByIDs(rIDs, nil)
	if err != nil {
		return false, err
	}

	j := helpers.JoinUintArrays(upIDs, pIDs)

	for _, p := range permissions.IDs() {
		if !helpers.InArray(p, j) {
			return false, err
		}
	}

	return true, err
}

// UserHasAnyPermissions */
// Does the user have any of the given permissions? (including the permissions of the roles)
func (s *Permify) UserHasAnyPermissions(userID uint, p interface{}) (b bool, err error) {

	var permissions collections.Permission
	permissions, err = s.GetPermissions(p)
	if err != nil {
		return false, err
	}

	var hasDirect bool
	hasDirect, err = s.UserRepository.HasAnyDirectPermissions(userID, permissions)
	if hasDirect {
		return true, nil
	}

	var roleIDs []uint
	roleIDs, _, err = s.RoleRepository.GetRoleIDsOfUser(userID, nil)

	var hasAny bool
	hasAny, err = s.RoleHasAnyPermissions(roleIDs, p)
	if hasAny {
		return true, nil
	}

	return false, err
}