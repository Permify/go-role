package repositories

import (
	"database/sql"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Permify/permify-gorm/collections"
	"github.com/Permify/permify-gorm/models"
	"github.com/Permify/permify-gorm/models/pivot"
)

var _ = Describe("User Repository", func() {
	var repository *UserRepository
	var mock sqlmock.Sqlmock

	BeforeEach(func() {
		var db *sql.DB
		var err error

		db, mock, err = sqlmock.New()
		Expect(err).ShouldNot(HaveOccurred())

		var gormDb *gorm.DB
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gormDb, err = gorm.Open(dialector, &gorm.Config{})
		Expect(err).ShouldNot(HaveOccurred())

		repository = &UserRepository{Database: gormDb}
	})

	AfterEach(func() {
		err := mock.ExpectationsWereMet()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("Has Role", func() {
		It("found", func() {
			userRoles := pivot.UserRoles{
				UserID: 1,
				RoleID: 1,
			}

			const query = `SELECT count(*) FROM "user_roles" WHERE user_roles.user_id = $1 AND user_roles.role_id = $2`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(userRoles.UserID, userRoles.RoleID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(1))

			value, err := repository.HasRole(1, models.Role{ID: 1})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(true))
		})

		It("not found", func() {
			const query = `SELECT count(*) FROM "user_roles" WHERE user_roles.user_id = $1 AND user_roles.role_id = $2`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(1, 1).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(0))

			value, err := repository.HasRole(1, models.Role{ID: 1})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(false))
		})
	})

	Context("Has All Roles", func() {
		It("found", func() {
			userRoles1 := pivot.UserRoles{
				UserID: 1,
				RoleID: 1,
			}

			userRoles2 := pivot.UserRoles{
				UserID: 1,
				RoleID: 2,
			}

			const query = `SELECT count(*) FROM "user_roles" WHERE user_roles.user_id = $1 AND user_roles.role_id IN ($2,$3)`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(userRoles1.UserID, userRoles1.RoleID, userRoles2.RoleID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(2))

			value, err := repository.HasAllRoles(uint(1), collections.Role([]models.Role{{ID: 1}, {ID: 2}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(true))
		})

		It("not found", func() {
			userRoles1 := pivot.UserRoles{
				UserID: 1,
				RoleID: 1,
			}

			userRoles2 := pivot.UserRoles{
				UserID: 1,
				RoleID: 2,
			}

			const query = `SELECT count(*) FROM "user_roles" WHERE user_roles.user_id = $1 AND user_roles.role_id IN ($2,$3)`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(userRoles1.UserID, userRoles1.RoleID, userRoles2.RoleID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(1))

			value, err := repository.HasAllRoles(uint(1), collections.Role([]models.Role{{ID: 1}, {ID: 2}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(false))
		})
	})

	Context("Has Any Roles", func() {
		It("found", func() {
			userRoles1 := pivot.UserRoles{
				UserID: 1,
				RoleID: 1,
			}

			const query = `SELECT count(*) FROM "user_roles" WHERE user_roles.user_id = $1 AND user_roles.role_id IN ($2,$3)`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(userRoles1.UserID, userRoles1.RoleID, 2).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(1))

			value, err := repository.HasAnyRoles(uint(1), collections.Role([]models.Role{{ID: 1}, {ID: 2}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(true))
		})

		It("not found", func() {
			const query = `SELECT count(*) FROM "user_roles" WHERE user_roles.user_id = $1 AND user_roles.role_id IN ($2,$3)`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(1, 1, 2).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(0))

			value, err := repository.HasAllRoles(uint(1), collections.Role([]models.Role{{ID: 1}, {ID: 2}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(false))
		})
	})

	Context("Has Direct Permission", func() {
		It("found", func() {
			userPermissions := pivot.UserPermissions{
				UserID:       1,
				PermissionID: 1,
			}

			const query = `SELECT count(*) FROM "user_permissions" WHERE user_permissions.user_id = $1 AND user_permissions.permission_id = $2`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(userPermissions.UserID, userPermissions.PermissionID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(1))

			value, err := repository.HasDirectPermission(1, models.Permission{ID: 1})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(true))
		})

		It("not found", func() {
			const query = `SELECT count(*) FROM "user_permissions" WHERE user_permissions.user_id = $1 AND user_permissions.permission_id = $2`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(1, 1).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(0))

			value, err := repository.HasDirectPermission(1, models.Permission{ID: 1})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(false))
		})
	})

	Context("Has All Direct Permissions", func() {
		It("found", func() {
			userPermissions1 := pivot.UserPermissions{
				UserID:       1,
				PermissionID: 1,
			}

			userPermissions2 := pivot.UserPermissions{
				UserID:       1,
				PermissionID: 2,
			}

			const query = `SELECT count(*) FROM "user_permissions" WHERE user_permissions.user_id = $1 AND user_permissions.permission_id IN ($2,$3)`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(userPermissions1.UserID, userPermissions1.PermissionID, userPermissions2.PermissionID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(2))

			value, err := repository.HasAllDirectPermissions(uint(1), collections.Permission([]models.Permission{{ID: 1}, {ID: 2}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(true))
		})

		It("not found", func() {
			userPermissions1 := pivot.UserPermissions{
				UserID:       1,
				PermissionID: 1,
			}

			userPermissions2 := pivot.UserPermissions{
				UserID:       1,
				PermissionID: 2,
			}

			const query = `SELECT count(*) FROM "user_permissions" WHERE user_permissions.user_id = $1 AND user_permissions.permission_id IN ($2,$3)`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(userPermissions1.UserID, userPermissions1.PermissionID, userPermissions2.PermissionID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(1))

			value, err := repository.HasAllDirectPermissions(uint(1), collections.Permission([]models.Permission{{ID: 1}, {ID: 2}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(false))
		})
	})

	Context("Has Any Direct Permissions", func() {
		It("found", func() {
			userPermissions := pivot.UserPermissions{
				UserID:       1,
				PermissionID: 1,
			}

			const query = `SELECT count(*) FROM "user_permissions" WHERE user_permissions.user_id = $1 AND user_permissions.permission_id IN ($2,$3)`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(userPermissions.UserID, userPermissions.PermissionID, 2).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(1))

			value, err := repository.HasAnyDirectPermissions(uint(1), collections.Permission([]models.Permission{{ID: 1}, {ID: 2}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(true))
		})

		It("not found", func() {
			const query = `SELECT count(*) FROM "user_permissions" WHERE user_permissions.user_id = $1 AND user_permissions.permission_id IN ($2,$3)`

			mock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(1, 1, 2).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(0))

			value, err := repository.HasAnyDirectPermissions(uint(1), collections.Permission([]models.Permission{{ID: 1}, {ID: 2}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(false))
		})
	})
})
