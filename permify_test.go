package permify_gorm

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Permify/permify-gorm/collections"
	"github.com/Permify/permify-gorm/models"
	"github.com/Permify/permify-gorm/options"
	"github.com/Permify/permify-gorm/repositories/mocks"
	"github.com/Permify/permify-gorm/repositories/scopes"
	"github.com/Permify/permify-gorm/utils"
)

func TestPermify(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Permify")
}

var _ = Describe("Permify", func() {
	var permify *Permify

	Context("Get Role", func() {
		It("By ID", func() {
			roleRepository := new(mocks.RoleRepository)
			r := models.Role{
				ID: 1,
			}
			roleRepository.On("GetRoleByID", uint(1)).Return(r, nil)
			permify = &Permify{
				RoleRepository: roleRepository,
			}
			actualResult, err := permify.GetRole(1, false)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(actualResult))
		})

		It("By Name", func() {
			roleRepository := new(mocks.RoleRepository)
			r := models.Role{
				Name:      "test role",
				GuardName: "test-role",
			}
			roleRepository.On("GetRoleByGuardName", "test-role").Return(r, nil)
			permify = &Permify{
				RoleRepository: roleRepository,
			}
			actualResult, err := permify.GetRole("test-role", false)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(actualResult))
		})
	})

	Context("Get Roles", func() {
		It("By IDs", func() {
			roleRepository := new(mocks.RoleRepository)
			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}
			roleRepository.On("GetRoles", []uint{1, 2}).Return(collections.Role(r), nil)
			permify = &Permify{
				RoleRepository: roleRepository,
			}
			actualResult, err := permify.GetRoles([]uint{1, 2}, false)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(actualResult.Origin()))
		})

		It("By Names", func() {
			roleRepository := new(mocks.RoleRepository)
			r := []models.Role{
				{
					Name:      "test role",
					GuardName: "test-role",
				},
				{
					Name:      "test role 2",
					GuardName: "test-role-2",
				},
			}
			roleRepository.On("GetRolesByGuardNames", []string{"test-role", "test-role-2"}).Return(collections.Role(r), nil)
			permify = &Permify{
				RoleRepository: roleRepository,
			}
			actualResult, err := permify.GetRoles([]string{"test role", "test role 2"}, false)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(actualResult.Origin()))
		})
	})

	Context("Get All Roles", func() {
		It("No Pagination", func() {
			roleRepository := new(mocks.RoleRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			roleRepository.On("GetRoleIDs", nil).Return([]uint{1, 2}, int64(2), nil)
			roleRepository.On("GetRoles", []uint{1, 2}).Return(collections.Role(r), nil)

			permify = &Permify{
				RoleRepository: roleRepository,
			}

			actualResult, _, err := permify.GetAllRoles(options.RoleOption{WithPermissions: false})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(actualResult.Origin()))
		})

		It("With Pagination", func() {
			roleRepository := new(mocks.RoleRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			roleRepository.On("GetRoleIDs", &scopes.GormPagination{
				Pagination: &utils.Pagination{
					Page:  1,
					Limit: 1,
				},
			}).Return([]uint{1}, int64(1), nil)
			roleRepository.On("GetRoles", []uint{1}).Return(collections.Role(r), nil)

			permify = &Permify{
				RoleRepository: roleRepository,
			}

			actualResult, _, err := permify.GetAllRoles(options.RoleOption{WithPermissions: false, Pagination: &utils.Pagination{Page: 1, Limit: 1}})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(actualResult.Origin()))
		})
	})

	Context("Get Roles of User", func() {
		It("No Pagination", func() {
			roleRepository := new(mocks.RoleRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			roleRepository.On("GetRoleIDsOfUser", uint(1), nil).Return([]uint{1, 2}, int64(2), nil)
			roleRepository.On("GetRoles", []uint{1, 2}).Return(collections.Role(r), nil)

			permify = &Permify{
				RoleRepository: roleRepository,
			}

			actualResult, _, err := permify.GetRolesOfUser(1, options.RoleOption{WithPermissions: false})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(actualResult.Origin()))
		})

		It("With Pagination", func() {
			roleRepository := new(mocks.RoleRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			roleRepository.On("GetRoleIDsOfUser", uint(1), &scopes.GormPagination{
				Pagination: &utils.Pagination{
					Page:  1,
					Limit: 1,
				},
			}).Return([]uint{1}, int64(1), nil)
			roleRepository.On("GetRoles", []uint{1}).Return(collections.Role(r), nil)

			permify = &Permify{
				RoleRepository: roleRepository,
			}

			actualResult, _, err := permify.GetRolesOfUser(1, options.RoleOption{WithPermissions: false, Pagination: &utils.Pagination{Page: 1, Limit: 1}})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(actualResult.Origin()))
		})
	})

	Context("Create Role", func() {
		It("Success", func() {
			roleRepository := new(mocks.RoleRepository)

			r := models.Role{
				Name:      "test",
				GuardName: "test",
			}

			roleRepository.On("FirstOrCreate", &r).Return(nil)

			permify = &Permify{
				RoleRepository: roleRepository,
			}

			Expect(permify.CreateRole(r.Name, r.Description)).ShouldNot(HaveOccurred())
		})
	})

	Context("Delete Role", func() {
		It("By ID", func() {
			roleRepository := new(mocks.RoleRepository)

			r := models.Role{
				ID:        1,
				Name:      "test",
				GuardName: "test",
			}

			roleRepository.On("GetRoleByID", uint(1)).Return(r, nil)
			roleRepository.On("Delete", &r).Return(nil)

			permify = &Permify{
				RoleRepository: roleRepository,
			}

			Expect(permify.DeleteRole(uint(1))).ShouldNot(HaveOccurred())
		})

		It("By Name", func() {
			roleRepository := new(mocks.RoleRepository)

			r := models.Role{
				Name:      "test",
				GuardName: "test",
			}

			roleRepository.On("GetRoleByGuardName", "test").Return(r, nil)
			roleRepository.On("Delete", &r).Return(nil)

			permify = &Permify{
				RoleRepository: roleRepository,
			}

			Expect(permify.DeleteRole("test")).ShouldNot(HaveOccurred())
		})
	})

	Context("Add Permissions to Role", func() {
		It("By IDs", func() {
			roleRepository := new(mocks.RoleRepository)
			permissionRepository := new(mocks.PermissionRepository)

			r := models.Role{
				Name:      "test role",
				GuardName: "test-role",
			}

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)
			roleRepository.On("GetRoleByGuardName", "test-role").Return(r, nil)
			roleRepository.On("AddPermissions", &r, collections.Permission(p)).Return(nil)
			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}
			Expect(permify.AddPermissionsToRole("test role", collections.Permission(p).IDs())).ShouldNot(HaveOccurred())
		})

		It("By Names", func() {
			roleRepository := new(mocks.RoleRepository)
			permissionRepository := new(mocks.PermissionRepository)

			r := models.Role{
				ID: 1,
			}

			p := []models.Permission{
				{
					Name:      "permission 1",
					GuardName: "permission-1",
				},
				{
					Name:      "permission 2",
					GuardName: "permission-2",
				},
			}

			permissionRepository.On("GetPermissionsByGuardNames", []string{"permission-1", "permission-2"}).Return(collections.Permission(p), nil)
			roleRepository.On("GetRoleByID", uint(1)).Return(r, nil)
			roleRepository.On("AddPermissions", &r, collections.Permission(p)).Return(nil)
			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}
			Expect(permify.AddPermissionsToRole(uint(1), collections.Permission(p).Names())).ShouldNot(HaveOccurred())
		})
	})

	Context("Replace Permissions to Role", func() {
		It("By IDs", func() {
			roleRepository := new(mocks.RoleRepository)
			permissionRepository := new(mocks.PermissionRepository)
			r := models.Role{
				Name:      "test role",
				GuardName: "test-role",
			}

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)
			roleRepository.On("GetRoleByGuardName", "test-role").Return(r, nil)
			roleRepository.On("ReplacePermissions", &r, collections.Permission(p)).Return(nil)
			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}
			Expect(permify.ReplacePermissionsToRole("test role", collections.Permission(p).IDs())).ShouldNot(HaveOccurred())
		})

		It("By Names", func() {
			roleRepository := new(mocks.RoleRepository)
			permissionRepository := new(mocks.PermissionRepository)

			r := models.Role{
				ID: 1,
			}

			p := []models.Permission{
				{
					Name:      "permission 1",
					GuardName: "permission-1",
				},
				{
					Name:      "permission 2",
					GuardName: "permission-2",
				},
			}

			permissionRepository.On("GetPermissionsByGuardNames", []string{"permission-1", "permission-2"}).Return(collections.Permission(p), nil)
			roleRepository.On("GetRoleByID", uint(1)).Return(r, nil)
			roleRepository.On("ReplacePermissions", &r, collections.Permission(p)).Return(nil)
			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}
			Expect(permify.ReplacePermissionsToRole(uint(1), collections.Permission(p).Names())).ShouldNot(HaveOccurred())
		})

		It("Clear", func() {
			roleRepository := new(mocks.RoleRepository)
			permissionRepository := new(mocks.PermissionRepository)
			r := models.Role{
				Name:      "test role",
				GuardName: "test-role",
			}

			permissionRepository.On("GetPermissions", []uint{}).Return(collections.Permission{}, nil)
			roleRepository.On("GetRoleByGuardName", "test-role").Return(r, nil)
			roleRepository.On("ClearPermissions", &r).Return(nil)
			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}
			Expect(permify.ReplacePermissionsToRole("test role", []uint{})).ShouldNot(HaveOccurred())
		})
	})

	Context("Remove Permissions from Role", func() {
		It("By IDs", func() {
			roleRepository := new(mocks.RoleRepository)
			permissionRepository := new(mocks.PermissionRepository)
			r := models.Role{
				Name:      "test role",
				GuardName: "test-role",
			}

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)
			roleRepository.On("GetRoleByGuardName", "test-role").Return(r, nil)
			roleRepository.On("RemovePermissions", &r, collections.Permission(p)).Return(nil)
			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}
			Expect(permify.RemovePermissionsFromRole("test role", collections.Permission(p).IDs())).ShouldNot(HaveOccurred())
		})

		It("By Names", func() {
			roleRepository := new(mocks.RoleRepository)
			permissionRepository := new(mocks.PermissionRepository)

			r := models.Role{
				ID: 1,
			}

			p := []models.Permission{
				{
					Name:      "permission 1",
					GuardName: "permission-1",
				},
				{
					Name:      "permission 2",
					GuardName: "permission-2",
				},
			}

			permissionRepository.On("GetPermissionsByGuardNames", []string{"permission-1", "permission-2"}).Return(collections.Permission(p), nil)
			roleRepository.On("GetRoleByID", uint(1)).Return(r, nil)
			roleRepository.On("RemovePermissions", &r, collections.Permission(p)).Return(nil)
			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}
			Expect(permify.RemovePermissionsFromRole(uint(1), collections.Permission(p).Names())).ShouldNot(HaveOccurred())
		})
	})

	Context("Get Permission", func() {
		It("By ID", func() {
			permissionRepository := new(mocks.PermissionRepository)
			r := models.Permission{
				ID: 1,
			}
			permissionRepository.On("GetPermissionByID", uint(1)).Return(r, nil)
			permify = &Permify{
				PermissionRepository: permissionRepository,
			}
			actualResult, err := permify.GetPermission(1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(actualResult))
		})

		It("By Name", func() {
			permissionRepository := new(mocks.PermissionRepository)
			r := models.Permission{
				Name:      "test permission",
				GuardName: "test-permission",
			}
			permissionRepository.On("GetPermissionByGuardName", "test-permission").Return(r, nil)
			permify = &Permify{
				PermissionRepository: permissionRepository,
			}
			actualResult, err := permify.GetPermission("test-permission")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(actualResult))
		})
	})

	Context("Get Permissions", func() {
		It("By IDs", func() {
			permissionRepository := new(mocks.PermissionRepository)
			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}
			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)
			permify = &Permify{
				PermissionRepository: permissionRepository,
			}
			actualResult, err := permify.GetPermissions([]uint{1, 2})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(Equal(actualResult.Origin()))
		})

		It("By Names", func() {
			permissionRepository := new(mocks.PermissionRepository)
			p := []models.Permission{
				{
					Name:      "test permission",
					GuardName: "test-permission",
				},
				{
					Name:      "test permission 2",
					GuardName: "test-permission-2",
				},
			}
			permissionRepository.On("GetPermissionsByGuardNames", []string{"test-permission", "test-permission-2"}).Return(collections.Permission(p), nil)
			permify = &Permify{
				PermissionRepository: permissionRepository,
			}
			actualResult, err := permify.GetPermissions([]string{"test permission", "test permission 2"})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(Equal(actualResult.Origin()))
		})
	})

	Context("Get All Permissions", func() {
		It("No Pagination", func() {
			permissionRepository := new(mocks.PermissionRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissionIDs", nil).Return([]uint{1, 2}, int64(2), nil)
			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
			}

			actualResult, _, err := permify.GetAllPermissions(options.PermissionOption{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(Equal(actualResult.Origin()))
		})

		It("With Pagination", func() {
			permissionRepository := new(mocks.PermissionRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissionIDs", &scopes.GormPagination{
				Pagination: &utils.Pagination{
					Page:  1,
					Limit: 1,
				},
			}).Return([]uint{1}, int64(1), nil)
			permissionRepository.On("GetPermissions", []uint{1}).Return(collections.Permission(p), nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
			}

			actualResult, _, err := permify.GetAllPermissions(options.PermissionOption{Pagination: &utils.Pagination{Page: 1, Limit: 1}})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(Equal(actualResult.Origin()))
		})
	})

	Context("Get Direct Permissions of User", func() {
		It("No Pagination", func() {
			permissionRepository := new(mocks.PermissionRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetDirectPermissionIDsOfUserByID", uint(1), nil).Return([]uint{1, 2}, int64(2), nil)
			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
			}

			actualResult, _, err := permify.GetDirectPermissionsOfUser(1, options.PermissionOption{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(Equal(actualResult.Origin()))
		})

		It("With Pagination", func() {
			permissionRepository := new(mocks.PermissionRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetDirectPermissionIDsOfUserByID", uint(1), &scopes.GormPagination{
				Pagination: &utils.Pagination{
					Page:  1,
					Limit: 1,
				},
			}).Return([]uint{1}, int64(1), nil)
			permissionRepository.On("GetPermissions", []uint{1}).Return(collections.Permission(p), nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
			}

			actualResult, _, err := permify.GetDirectPermissionsOfUser(1, options.PermissionOption{Pagination: &utils.Pagination{Page: 1, Limit: 1}})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(Equal(actualResult.Origin()))
		})
	})

	Context("Get Permissions of Roles", func() {
		It("No Pagination", func() {
			permissionRepository := new(mocks.PermissionRepository)
			roleRepository := new(mocks.RoleRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissionIDsOfRolesByIDs", []uint{1, 2}, nil).Return([]uint{1, 2}, int64(2), nil)
			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)
			roleRepository.On("GetRoles", []uint{1, 2}).Return(collections.Role(r), nil)

			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}

			actualResult, _, err := permify.GetPermissionsOfRoles(collections.Role(r).IDs(), options.PermissionOption{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(Equal(actualResult.Origin()))
		})

		It("With Pagination", func() {
			permissionRepository := new(mocks.PermissionRepository)
			roleRepository := new(mocks.RoleRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissionIDsOfRolesByIDs", []uint{1, 2}, &scopes.GormPagination{
				Pagination: &utils.Pagination{
					Page:  1,
					Limit: 1,
				},
			}).Return([]uint{1}, int64(1), nil)
			permissionRepository.On("GetPermissions", []uint{1}).Return(collections.Permission(p), nil)
			roleRepository.On("GetRoles", []uint{1, 2}).Return(collections.Role(r), nil)

			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}

			actualResult, _, err := permify.GetPermissionsOfRoles(collections.Role(r).IDs(), options.PermissionOption{Pagination: &utils.Pagination{Page: 1, Limit: 1}})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(Equal(actualResult.Origin()))
		})
	})

	Context("Get All Permissions of User", func() {
		It("No Pagination", func() {
			permissionRepository := new(mocks.PermissionRepository)
			roleRepository := new(mocks.RoleRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			roleRepository.On("GetRoleIDsOfUser", uint(1), nil).Return([]uint{1, 2}, int64(2), nil)

			permissionRepository.On("GetDirectPermissionIDsOfUserByID", uint(1), nil).Return([]uint{1}, int64(1), nil)
			permissionRepository.On("GetPermissionIDsOfRolesByIDs", []uint{1, 2}, nil).Return([]uint{1, 2}, int64(2), nil)
			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)

			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}

			actualResult, err := permify.GetAllPermissionsOfUser(uint(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(p).Should(Equal(actualResult.Origin()))
		})
	})

	Context("Create Permission", func() {
		It("Success", func() {
			permissionRepository := new(mocks.PermissionRepository)

			p := models.Permission{
				Name:      "test",
				GuardName: "test",
			}

			permissionRepository.On("FirstOrCreate", &p).Return(nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
			}

			Expect(permify.CreatePermission(p.Name, p.Description)).ShouldNot(HaveOccurred())
		})
	})

	Context("Delete Permission", func() {
		It("By ID", func() {
			permissionRepository := new(mocks.PermissionRepository)

			r := models.Permission{
				ID:        1,
				Name:      "test",
				GuardName: "test",
			}

			permissionRepository.On("GetPermissionByID", uint(1)).Return(r, nil)
			permissionRepository.On("Delete", &r).Return(nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
			}

			Expect(permify.DeletePermission(uint(1))).ShouldNot(HaveOccurred())
		})

		It("By Name", func() {
			permissionRepository := new(mocks.PermissionRepository)

			p := models.Permission{
				Name:      "test",
				GuardName: "test",
			}

			permissionRepository.On("GetPermissionByGuardName", "test").Return(p, nil)
			permissionRepository.On("Delete", &p).Return(nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
			}

			Expect(permify.DeletePermission("test")).ShouldNot(HaveOccurred())
		})
	})

	Context("Add Permissions to User", func() {
		It("By IDs", func() {
			permissionRepository := new(mocks.PermissionRepository)
			userRepository := new(mocks.UserRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			userRepository.On("AddPermissions", uint(1), collections.Permission(p)).Return(nil)
			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)

			permify = &Permify{
				UserRepository:       userRepository,
				PermissionRepository: permissionRepository,
			}

			Expect(permify.AddPermissionsToUser(uint(1), collections.Permission(p).IDs())).ShouldNot(HaveOccurred())
		})

		It("By Names", func() {
			permissionRepository := new(mocks.PermissionRepository)
			userRepository := new(mocks.UserRepository)

			p := []models.Permission{
				{
					Name:      "permission 1",
					GuardName: "permission-1",
				},
				{
					Name:      "permission 2",
					GuardName: "permission-2",
				},
			}

			userRepository.On("AddPermissions", uint(1), collections.Permission(p)).Return(nil)
			permissionRepository.On("GetPermissionsByGuardNames", []string{"permission-1", "permission-2"}).Return(collections.Permission(p), nil)

			permify = &Permify{
				UserRepository:       userRepository,
				PermissionRepository: permissionRepository,
			}

			Expect(permify.AddPermissionsToUser(uint(1), collections.Permission(p).Names())).ShouldNot(HaveOccurred())
		})
	})

	Context("Replace Permissions to User", func() {
		It("By IDs", func() {
			permissionRepository := new(mocks.PermissionRepository)
			userRepository := new(mocks.UserRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			userRepository.On("ReplacePermissions", uint(1), collections.Permission(p)).Return(nil)
			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)

			permify = &Permify{
				UserRepository:       userRepository,
				PermissionRepository: permissionRepository,
			}

			Expect(permify.ReplacePermissionsToUser(uint(1), collections.Permission(p).IDs())).ShouldNot(HaveOccurred())
		})

		It("By Names", func() {
			permissionRepository := new(mocks.PermissionRepository)
			userRepository := new(mocks.UserRepository)

			p := []models.Permission{
				{
					Name:      "permission 1",
					GuardName: "permission-1",
				},
				{
					Name:      "permission 2",
					GuardName: "permission-2",
				},
			}

			userRepository.On("ReplacePermissions", uint(1), collections.Permission(p)).Return(nil)
			permissionRepository.On("GetPermissionsByGuardNames", []string{"permission-1", "permission-2"}).Return(collections.Permission(p), nil)

			permify = &Permify{
				UserRepository:       userRepository,
				PermissionRepository: permissionRepository,
			}

			Expect(permify.ReplacePermissionsToUser(uint(1), collections.Permission(p).Names())).ShouldNot(HaveOccurred())
		})

		It("Clear", func() {
			userRepository := new(mocks.UserRepository)
			permissionRepository := new(mocks.PermissionRepository)

			permissionRepository.On("GetPermissions", []uint{}).Return(collections.Permission{}, nil)
			userRepository.On("ClearPermissions", uint(1)).Return(nil)

			permify = &Permify{
				UserRepository:       userRepository,
				PermissionRepository: permissionRepository,
			}

			Expect(permify.ReplacePermissionsToUser(uint(1), []uint{})).ShouldNot(HaveOccurred())
		})
	})

	Context("Remove Permissions from User", func() {
		It("By IDs", func() {
			permissionRepository := new(mocks.PermissionRepository)
			userRepository := new(mocks.UserRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			userRepository.On("RemovePermissions", uint(1), collections.Permission(p)).Return(nil)
			permissionRepository.On("GetPermissions", []uint{1, 2}).Return(collections.Permission(p), nil)

			permify = &Permify{
				UserRepository:       userRepository,
				PermissionRepository: permissionRepository,
			}

			Expect(permify.RemovePermissionsFromUser(uint(1), collections.Permission(p).IDs())).ShouldNot(HaveOccurred())
		})

		It("By Names", func() {
			permissionRepository := new(mocks.PermissionRepository)
			userRepository := new(mocks.UserRepository)

			p := []models.Permission{
				{
					Name:      "permission 1",
					GuardName: "permission-1",
				},
				{
					Name:      "permission 2",
					GuardName: "permission-2",
				},
			}

			userRepository.On("RemovePermissions", uint(1), collections.Permission(p)).Return(nil)
			permissionRepository.On("GetPermissionsByGuardNames", []string{"permission-1", "permission-2"}).Return(collections.Permission(p), nil)

			permify = &Permify{
				UserRepository:       userRepository,
				PermissionRepository: permissionRepository,
			}

			Expect(permify.RemovePermissionsFromUser(uint(1), collections.Permission(p).Names())).ShouldNot(HaveOccurred())
		})
	})

	Context("Add Roles to User", func() {
		It("By IDs", func() {
			roleRepository := new(mocks.RoleRepository)
			userRepository := new(mocks.UserRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			roleRepository.On("GetRoles", collections.Role(r).IDs()).Return(collections.Role(r), nil)
			userRepository.On("AddRoles", uint(1), collections.Role(r)).Return(nil)

			permify = &Permify{
				UserRepository: userRepository,
				RoleRepository: roleRepository,
			}

			Expect(permify.AddRolesToUser(uint(1), collections.Role(r).IDs())).ShouldNot(HaveOccurred())
		})

		It("By Names", func() {
			roleRepository := new(mocks.RoleRepository)
			userRepository := new(mocks.UserRepository)

			r := []models.Role{
				{
					Name:      "role 1",
					GuardName: "role-1",
				},
				{
					Name:      "role 2",
					GuardName: "role-2",
				},
			}

			roleRepository.On("GetRolesByGuardNames", collections.Role(r).GuardNames()).Return(collections.Role(r), nil)
			userRepository.On("AddRoles", uint(1), collections.Role(r)).Return(nil)

			permify = &Permify{
				UserRepository: userRepository,
				RoleRepository: roleRepository,
			}

			Expect(permify.AddRolesToUser(uint(1), collections.Role(r).Names())).ShouldNot(HaveOccurred())
		})
	})

	Context("Replace Roles to User", func() {
		It("By IDs", func() {
			roleRepository := new(mocks.RoleRepository)
			userRepository := new(mocks.UserRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			roleRepository.On("GetRoles", collections.Role(r).IDs()).Return(collections.Role(r), nil)
			userRepository.On("ReplaceRoles", uint(1), collections.Role(r)).Return(nil)

			permify = &Permify{
				UserRepository: userRepository,
				RoleRepository: roleRepository,
			}

			Expect(permify.ReplaceRolesToUser(uint(1), collections.Role(r).IDs())).ShouldNot(HaveOccurred())
		})

		It("By Names", func() {
			roleRepository := new(mocks.RoleRepository)
			userRepository := new(mocks.UserRepository)

			r := []models.Role{
				{
					Name:      "role 1",
					GuardName: "role-1",
				},
				{
					Name:      "role 2",
					GuardName: "role-2",
				},
			}

			roleRepository.On("GetRolesByGuardNames", collections.Role(r).GuardNames()).Return(collections.Role(r), nil)
			userRepository.On("ReplaceRoles", uint(1), collections.Role(r)).Return(nil)

			permify = &Permify{
				UserRepository: userRepository,
				RoleRepository: roleRepository,
			}

			Expect(permify.ReplaceRolesToUser(uint(1), collections.Role(r).Names())).ShouldNot(HaveOccurred())
		})

		It("Clear", func() {
			roleRepository := new(mocks.RoleRepository)
			userRepository := new(mocks.UserRepository)

			roleRepository.On("GetRoles", []uint{}).Return(collections.Role{}, nil)
			userRepository.On("ClearRoles", uint(1)).Return(nil)

			permify = &Permify{
				UserRepository: userRepository,
				RoleRepository: roleRepository,
			}

			Expect(permify.ReplaceRolesToUser(uint(1), []uint{})).ShouldNot(HaveOccurred())
		})
	})

	// Controls

	Context("Role Has Permission", func() {
		It("Success", func() {
			roleRepository := new(mocks.RoleRepository)
			permissionRepository := new(mocks.PermissionRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			p := models.Permission{
				ID:        1,
				Name:      "permission 1",
				GuardName: "permission-1",
			}

			roleRepository.On("GetRoles", collections.Role(r).IDs()).Return(collections.Role(r), nil)
			permissionRepository.On("GetPermissionByID", p.ID).Return(p, nil)
			roleRepository.On("HasPermission", collections.Role(r), p).Return(true, nil)

			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}

			actualResult, err := permify.RoleHasPermission(collections.Role(r).IDs(), p.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("Role Has All Permission", func() {
		It("Success", func() {
			roleRepository := new(mocks.RoleRepository)
			permissionRepository := new(mocks.PermissionRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			p := []models.Permission{
				{
					ID:        1,
					Name:      "permission 1",
					GuardName: "permission-1",
				},
				{
					ID:        2,
					Name:      "permission 2",
					GuardName: "permission-2",
				},
			}

			roleRepository.On("GetRoles", collections.Role(r).IDs()).Return(collections.Role(r), nil)
			permissionRepository.On("GetPermissions", collections.Permission(p).IDs()).Return(collections.Permission(p), nil)
			roleRepository.On("HasAllPermissions", collections.Role(r), collections.Permission(p)).Return(true, nil)

			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}

			actualResult, err := permify.RoleHasAllPermissions(collections.Role(r).IDs(), collections.Role(r).IDs())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("Role Has Any Permission", func() {
		It("Success", func() {
			roleRepository := new(mocks.RoleRepository)
			permissionRepository := new(mocks.PermissionRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			p := []models.Permission{
				{
					ID:        1,
					Name:      "permission 1",
					GuardName: "permission-1",
				},
				{
					ID:        2,
					Name:      "permission 2",
					GuardName: "permission-2",
				},
			}

			roleRepository.On("GetRoles", collections.Role(r).IDs()).Return(collections.Role(r), nil)
			permissionRepository.On("GetPermissions", collections.Permission(p).IDs()).Return(collections.Permission(p), nil)
			roleRepository.On("HasAnyPermissions", collections.Role(r), collections.Permission(p)).Return(true, nil)

			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}

			actualResult, err := permify.RoleHasAnyPermissions(collections.Role(r).IDs(), collections.Role(r).IDs())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("User Has Role", func() {
		It("Success", func() {
			roleRepository := new(mocks.RoleRepository)
			userRepository := new(mocks.UserRepository)

			r := models.Role{
				ID: 1,
			}

			roleRepository.On("GetRoleByID", r.ID).Return(r, nil)
			userRepository.On("HasRole", uint(1), r).Return(true, nil)

			permify = &Permify{
				RoleRepository: roleRepository,
				UserRepository: userRepository,
			}

			actualResult, err := permify.UserHasRole(uint(1), r.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("User Has All Roles", func() {
		It("Success", func() {
			roleRepository := new(mocks.RoleRepository)
			userRepository := new(mocks.UserRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			roleRepository.On("GetRoles", collections.Role(r).IDs()).Return(collections.Role(r), nil)
			userRepository.On("HasAllRoles", uint(1), collections.Role(r)).Return(true, nil)

			permify = &Permify{
				RoleRepository: roleRepository,
				UserRepository: userRepository,
			}

			actualResult, err := permify.UserHasAllRoles(uint(1), collections.Role(r).IDs())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("User Has Any Roles", func() {
		It("Success", func() {
			roleRepository := new(mocks.RoleRepository)
			userRepository := new(mocks.UserRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			roleRepository.On("GetRoles", collections.Role(r).IDs()).Return(collections.Role(r), nil)
			userRepository.On("HasAnyRoles", uint(1), collections.Role(r)).Return(true, nil)

			permify = &Permify{
				RoleRepository: roleRepository,
				UserRepository: userRepository,
			}

			actualResult, err := permify.UserHasAnyRoles(uint(1), collections.Role(r).IDs())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("User Has Direct Permission", func() {
		It("Success", func() {
			permissionRepository := new(mocks.PermissionRepository)
			userRepository := new(mocks.UserRepository)

			p := models.Permission{
				ID: 1,
			}

			permissionRepository.On("GetPermissionByID", p.ID).Return(p, nil)
			userRepository.On("HasDirectPermission", uint(1), p).Return(true, nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
				UserRepository:       userRepository,
			}

			actualResult, err := permify.UserHasDirectPermission(uint(1), p.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("User Has All Direct Permissions", func() {
		It("Success", func() {
			permissionRepository := new(mocks.PermissionRepository)
			userRepository := new(mocks.UserRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissions", collections.Permission(p).IDs()).Return(collections.Permission(p), nil)
			userRepository.On("HasAllDirectPermissions", uint(1), collections.Permission(p)).Return(true, nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
				UserRepository:       userRepository,
			}

			actualResult, err := permify.UserHasAllDirectPermissions(uint(1), collections.Permission(p).IDs())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("User Has Any Direct Permissions", func() {
		It("Success", func() {
			permissionRepository := new(mocks.PermissionRepository)
			userRepository := new(mocks.UserRepository)

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissions", collections.Permission(p).IDs()).Return(collections.Permission(p), nil)
			userRepository.On("HasAnyDirectPermissions", uint(1), collections.Permission(p)).Return(true, nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
				UserRepository:       userRepository,
			}

			actualResult, err := permify.UserHasAnyDirectPermissions(uint(1), collections.Permission(p).IDs())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("User Has Permission", func() {
		It("Has Direct Permission Success", func() {
			permissionRepository := new(mocks.PermissionRepository)

			p := models.Permission{
				ID: 1,
			}

			permissionRepository.On("GetPermissionByID", p.ID).Return(p, nil)
			permissionRepository.On("GetDirectPermissionIDsOfUserByID", uint(1), nil).Return([]uint{1}, int64(1), nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
			}

			actualResult, err := permify.UserHasPermission(uint(1), p.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})

		It("Non Direct Permission Success", func() {
			permissionRepository := new(mocks.PermissionRepository)
			roleRepository := new(mocks.RoleRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			p := models.Permission{
				ID: 1,
			}

			permissionRepository.On("GetPermissionByID", p.ID).Return(p, nil)
			permissionRepository.On("GetDirectPermissionIDsOfUserByID", uint(1), nil).Return([]uint{3}, int64(1), nil)
			roleRepository.On("GetRoleIDsOfUser", uint(1), nil).Return(collections.Role(r).IDs(), int64(1), nil)
			permissionRepository.On("GetPermissionIDsOfRolesByIDs", collections.Role(r).IDs(), nil).Return([]uint{1}, int64(1), nil)

			permify = &Permify{
				PermissionRepository: permissionRepository,
				RoleRepository:       roleRepository,
			}

			actualResult, err := permify.UserHasPermission(uint(1), p.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("User Has All Permissions", func() {
		It("Success", func() {
			permissionRepository := new(mocks.PermissionRepository)
			roleRepository := new(mocks.RoleRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissions", collections.Permission(p).IDs()).Return(collections.Permission(p), nil)
			permissionRepository.On("GetDirectPermissionIDsOfUserByID", uint(1), nil).Return([]uint{1, 2}, int64(1), nil)
			roleRepository.On("GetRoleIDsOfUser", uint(1), nil).Return(collections.Role(r).IDs(), int64(1), nil)
			permissionRepository.On("GetPermissionIDsOfRolesByIDs", collections.Role(r).IDs(), nil).Return([]uint{1}, int64(1), nil)

			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}

			actualResult, err := permify.UserHasAllPermissions(uint(1), collections.Permission(p).IDs())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})

	Context("User Has Any Permissions", func() {
		It("Success", func() {
			permissionRepository := new(mocks.PermissionRepository)
			roleRepository := new(mocks.RoleRepository)

			r := []models.Role{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			p := []models.Permission{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			}

			permissionRepository.On("GetPermissions", collections.Permission(p).IDs()).Return(collections.Permission(p), nil)
			permissionRepository.On("GetDirectPermissionIDsOfUserByID", uint(1), nil).Return([]uint{1, 2}, int64(1), nil)
			roleRepository.On("GetRoleIDsOfUser", uint(1), nil).Return(collections.Role(r).IDs(), int64(1), nil)
			permissionRepository.On("GetPermissionIDsOfRolesByIDs", collections.Role(r).IDs(), nil).Return([]uint{1}, int64(1), nil)

			permify = &Permify{
				RoleRepository:       roleRepository,
				PermissionRepository: permissionRepository,
			}

			actualResult, err := permify.UserHasAnyPermissions(uint(1), collections.Permission(p).IDs())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(Equal(actualResult))
		})
	})
})
