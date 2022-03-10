
![permify-gorm](https://user-images.githubusercontent.com/39353278/157410086-42e02752-d5a9-4c64-bdc3-d3a203a247d7.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/Permify/permify-gorm)](https://goreportcard.com/report/github.com/Permify/permify-gorm)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Permify/permify-gorm)
![GitHub](https://img.shields.io/github/license/Permify/permify-gorm)

## Role Based Access Control (RBAC) for your go application

This package allows you to manage user permissions and roles in your database.


## 👇 Setup

Install

```shell
go get github.com/Permify/permify-gorm
```

Run All Tests

```shell
go test ./...
```

Get the database driver for gorm that you will be using

```shell
# mysql 
go get gorm.io/driver/mysql 
# or postgres
go get gorm.io/driver/postgres
# or sqlite
go get gorm.io/driver/sqlite
# or sqlserver
go get gorm.io/driver/sqlserver
# or clickhouse
go get gorm.io/driver/clickhouse
```

Import permify.

```go
import permify `github.com/Permify/permify-gorm`
```

Initialize the new Permify.

```go
// initialize the database. (you can use all gorm's supported databases)
db, _ := gorm.Open(mysql.Open("user:password@tcp(host:3306)/db?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})

// New initializer for Permify
// If migration is true, it generate all tables in the database if they don't exist.
permify, _ := permify.New(permify.Options{
	Migrate: true,
	DB: db,
})
```

## 🚲 Basic Usage

This package allows users to be associated with permissions and roles. Each role is associated with multiple permissions.

```go
// Create new role
// Name parameter is converted to guard name. example: senior $#% associate -> senior-associate.
// If a role with the same name has been created before, it will not create it again. (FirstOrCreate)
err := permify.CreateRole("admin", "role description")

// Create new permission
// Name parameter is converted to guard name. example: create $#% contact -> create-contact.
// If a permission with the same name has been created before, it will not create it again. (FirstOrCreate)
err := permify.CreatePermission("edit user details", "")
```

A permissions can be added to a role using this method in different ways:

```go
// first parameter is role id
err := p.AddPermissionsToRole(1, "edit user details")
// or
err := p.AddPermissionsToRole("admin", []string{"edit user details", "create contact"})
// or
err := p.AddPermissionsToRole("admin", []uint{1, 3})
```

Using these methods you can manage roles permissions removes and overwrites like the same above ways:

```go
// Overwrites the permissions of the role according to the permission names or ids.
err := permify.ReplacePermissionsToRole("admin", []string{"edit user details", "create contact"})

// remove permissions from role according to the permission names or ids.
err := permify.RemovePermissionsFromRole("admin", []string{"edit user details"})
```

Basic fetch queries:

```go
// Fetch all the roles. (with pagination option).
// If withPermissions is true, it will preload the permissions to the role.
// If pagination is nil, it returns without paging.
roles, totalCount, err := permify.GetAllRoles(options.RoleOption{
	WithPermissions: true,
	Pagination: &utils.Pagination{
		Page: 1,
		Limit: 1,
	},
})

// without paging.
roles, totalCount, err := permify.GetAllRoles(options.RoleOption{
    WithPermissions: false,
})

// The data returned is a collection of roles. Collections provides a fluent convenient wrapper for working with arrays of data.
fmt.Println(roles.IDs())
fmt.Println(roles.Names())
fmt.Println(roles.Permissions().Names())

// Fetch all permissions of the user that come with direct and roles.
permissions, _ := permify.GetAllPermissionsOfUser(1)

// Fetch all direct permissions of the user. (with pagination option)
permissions, totalCount, err := p.GetDirectPermissionsOfUser(1, options.PermissionOption{
    Pagination: &utils.Pagination{
        Page: 1,
        Limit: 10,
    },
})
```

Controls

```go
// Does the role or any of the roles have given permission?
can, err := permify.RoleHasPermission("admin", "edit user details")

// Does the role or roles have any of the given permissions?
can, err := permify.RoleHasAnyPermissions([]string{"admin", "manager"}, []string{"edit user details", "create contact"})

// Does the role or roles have all the given permissions?
can, err := permify.RoleHasAllPermissions("admin", []string{"edit user details", "create contact"})

// Does the user have the given permission? (including the permissions of the roles)
can, err := permify.UserHasPermission(1, "edit user details")

// Does the user have the given permission? (not including the permissions of the roles)
can, err := permify.UserHasDirectPermission(1, "edit user details")

// Does the user have any of the given permissions? (including the permissions of the roles)
can, err := permify.UserHasAnyPermissions(1, []uint{1, 2})

// Does the user have all the given roles?
can, err := permify.UserHasAllRoles(1, []string{"admin", "manager"})

// Does the user have any of the given roles?
can, err := permify.UserHasAnyRoles(1, []string{"admin", "manager"})
```


## 🚘 Using permissions via roles

### Adding Role

Add role or roles to user according to the role names or ids:

```go
// add one role to user
err := permify.AddRolesToUser(1, "admin")

// You can also add multiple roles at once
err := permify.AddRolesToUser(1, []string{"admin", "manager"})
// or
err := permify.AddRolesToUser(1, []uint{1,2})
```

Replace the roles of the user according to the role names or ids:

```go
// remove all user roles and add admin role
err := permify.ReplaceRolesToUser(1, "admin")

// You can also replace multiple roles at once
err := permify.ReplaceRolesToUser(1, []string{"admin", "manager"})
// or
err := permify.RemoveRolesFromUser(1, []uint{1,2})
```

Controls

```go

// does the user have the given role?
can, err := permify.UserHasRole(1, "admin")
```


## 🚤 Direct Permissions





## 🚀 Using your user model

You can create the relationships between the user and the role and permissions in this manner. In this way:

- You can manage user preloads
- You can create foreign key between users and pivot tables (user_roles, user_permissions).

```go
import models `github.com/Permify/permify-gorm/models`

type User struct {
    gorm.Model
    Name string

    // permify
    Roles []models.Role `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Permissions []models.Permission `gorm:"many2many:user_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
```

## ⁉️ Error Handling

### ErrRecordNotFound

You can use error handling in the same way as gorm. for example:

```go
// Check if returns RecordNotFound error
permission, err := permify.GetPermission(1)
if errors.Is(err, gorm.ErrRecordNotFound) {
	// record not found
}
```

### Errors

[Errors List](https://github.com/go-gorm/gorm/blob/master/errors.go)

<h2 align="left">:heart: Let's get connected:</h2>

<p align="left">
<a href="https://twitter.com/GetPermify">
  <img alt="guilyx | Twitter" width="50px" src="https://user-images.githubusercontent.com/43545812/144034996-602b144a-16e1-41cc-99e7-c6040b20dcaf.png"/>
</a>
<a href="https://www.linkedin.com/company/permifyco">
  <img alt="guilyx's LinkdeIN" width="50px" src="https://user-images.githubusercontent.com/43545812/144035037-0f415fc7-9f96-4517-a370-ccc6e78a714b.png" />
</a>
</p>


[comment]: <> (![permify-gorm-draw-sql]&#40;https://user-images.githubusercontent.com/39353278/157461050-0a146e7c-9ba7-4956-90a9-4720190a2c82.png&#41;)
