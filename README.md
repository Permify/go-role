
![permify-gorm](https://user-images.githubusercontent.com/39353278/157410086-42e02752-d5a9-4c64-bdc3-d3a203a247d7.png)

## Role Based Access Control (RBAC) for your go application

This package allows you to manage user permissions and roles in your database.

<br>

<details>
<summary>👇 Install</summary>

```shell
go get github.com/Permify/permify-gorm
```

Initialize the new Permify.

```go
// initialize the database. (you can use all gorm's supported databases)
db, _ := gorm.Open(mysql.Open("user:password@tcp(host:3306)/db?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})

// initialize the Permify
permify, _ := New(Options{
	Migrate: true,
	DB: db,
})
```

<br>

</details>

<br>

<details>
<summary>🚲 Basic Usage</summary>
<br>
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

A permission or permissions can be added to a role using this method in different ways:

```go
err := permify.AddPermissionsToRole(1, "edit user details")
// or
err := permify.AddPermissionsToRole("admin", []string{"edit user details", "create contact"})
// or
err := permify.AddPermissionsToRole("admin", []uint{1, 3})
```

</details>

<br>

<details>
<summary>🚘 Using permissions via roles</summary>
<br>

</details>

<br>

<details>
<summary>🚤 Direct Permissions</summary>
<br>

</details>

<br>

<details>
<summary>🚀 Using your user model</summary>
<br>

</details>

<br>

<h2 align="left">:heart: Let's get connected:</h2>

-----

<p align="left">
<a href="https://twitter.com/GetPermify">
  <img alt="guilyx | Twitter" width="50px" src="https://user-images.githubusercontent.com/43545812/144034996-602b144a-16e1-41cc-99e7-c6040b20dcaf.png"/>
</a>
<a href="https://www.linkedin.com/company/permifyco">
  <img alt="guilyx's LinkdeIN" width="50px" src="https://user-images.githubusercontent.com/43545812/144035037-0f415fc7-9f96-4517-a370-ccc6e78a714b.png" />
</a>
</p>


[comment]: <> (![permify-gorm-draw-sql]&#40;https://user-images.githubusercontent.com/39353278/157461050-0a146e7c-9ba7-4956-90a9-4720190a2c82.png&#41;)
