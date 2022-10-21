package repositories

import (
	"database/sql"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Permify/go-role/collections"
	"github.com/Permify/go-role/models"
)

var _ = Describe("Role Repository", func() {
	var repository *RoleRepository
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

		repository = &RoleRepository{Database: gormDb}
	})

	AfterEach(func() {
		err := mock.ExpectationsWereMet()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("Get Role By ID", func() {
		It("found", func() {
			role := models.Role{
				ID:        1,
				Name:      "admin",
				GuardName: "admin",
			}

			rows := sqlmock.NewRows([]string{"id", "name", "guard_name"}).
				AddRow(role.ID, role.Name, role.GuardName)

			const sqlSelectOne = `SELECT * FROM "roles" WHERE roles.id = $1 ORDER BY "roles"."id" LIMIT 1`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(role.ID).
				WillReturnRows(rows)

			db, err := repository.GetRoleByID(role.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(role))
		})

		It("not found", func() {
			// ignore sql match
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := repository.GetRoleByID(1)
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
	})

	Context("Get Role By Guard Name", func() {
		It("found", func() {
			role := models.Role{
				ID:        1,
				Name:      "admin",
				GuardName: "admin",
			}

			rows := sqlmock.NewRows([]string{"id", "name", "guard_name"}).
				AddRow(role.ID, role.Name, role.GuardName)

			const sqlSelectOne = `SELECT * FROM "roles" WHERE roles.guard_name = $1 ORDER BY "roles"."id" LIMIT 1`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(role.GuardName).
				WillReturnRows(rows)

			db, err := repository.GetRoleByGuardName(role.GuardName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(role))
		})

		It("not found", func() {
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := repository.GetRoleByGuardName("admin")
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
	})

	Context("Get Roles", func() {
		It("found", func() {
			roles := []models.Role{
				{
					ID:        1,
					Name:      "admin",
					GuardName: "admin",
				},
			}

			rows := sqlmock.NewRows([]string{"id", "name", "guard_name"}).
				AddRow(roles[0].ID, roles[0].Name, roles[0].GuardName)

			const sqlSelectOne = `SELECT * FROM "roles" WHERE roles.id IN ($1)`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(roles[0].ID).
				WillReturnRows(rows)

			value, err := repository.GetRoles([]uint{roles[0].ID})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value.Origin()).Should(Equal(roles))
		})
	})

	Context("Get Roles By Guard Names", func() {
		It("found", func() {
			roles := []models.Role{
				{
					ID:        1,
					Name:      "admin",
					GuardName: "admin",
				},
			}

			rows := sqlmock.NewRows([]string{"id", "name", "guard_name"}).
				AddRow(roles[0].ID, roles[0].Name, roles[0].GuardName)

			const sqlSelectOne = `SELECT * FROM "roles" WHERE roles.guard_name IN ($1)`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(roles[0].GuardName).
				WillReturnRows(rows)

			value, err := repository.GetRolesByGuardNames([]string{roles[0].GuardName})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value.Origin()).Should(Equal(roles))
		})
	})

	Context("Has Permission", func() {
		It("found", func() {
			const sqlSelectOne = `SELECT count(*) FROM "role_permissions" WHERE role_permissions.role_id IN ($1) AND role_permissions.permission_id = $2`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(1, 1).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(1))

			db, err := repository.HasPermission(collections.Role([]models.Role{{ID: 1}}), models.Permission{ID: 1})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(true))
		})

		It("not found", func() {
			const sqlSelectOne = `SELECT count(*) FROM "role_permissions" WHERE role_permissions.role_id IN ($1) AND role_permissions.permission_id = $2`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(1, 1).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(0))

			db, err := repository.HasPermission(collections.Role([]models.Role{{ID: 1}}), models.Permission{ID: 1})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(false))
		})
	})

	Context("Has All Permission", func() {
		It("found", func() {
			const sqlSelectOne = `SELECT count(*) FROM "role_permissions" WHERE role_permissions.role_id IN ($1) AND role_permissions.permission_id IN ($2)`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(1, 1).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(1))

			db, err := repository.HasAllPermissions(collections.Role([]models.Role{{ID: 1}}), collections.Permission([]models.Permission{{ID: 1}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(true))
		})

		It("not found", func() {
			const sqlSelectOne = `SELECT count(*) FROM "role_permissions" WHERE role_permissions.role_id IN ($1,$2) AND role_permissions.permission_id IN ($3)`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(1, 2, 1).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(1))

			db, err := repository.HasAllPermissions(collections.Role([]models.Role{{ID: 1}, {ID: 2}}), collections.Permission([]models.Permission{{ID: 1}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(false))
		})
	})

	Context("Has Any Permission", func() {
		It("found", func() {
			const sqlSelectOne = `SELECT count(*) FROM "role_permissions" WHERE role_permissions.role_id IN ($1) AND role_permissions.permission_id IN ($2)`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(1, 1).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(1))

			db, err := repository.HasAnyPermissions(collections.Role([]models.Role{{ID: 1}}), collections.Permission([]models.Permission{{ID: 1}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(true))
		})

		It("found", func() {
			const sqlSelectOne = `SELECT count(*) FROM "role_permissions" WHERE role_permissions.role_id IN ($1) AND role_permissions.permission_id IN ($2)`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(1, 1).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(0))

			db, err := repository.HasAnyPermissions(collections.Role([]models.Role{{ID: 1}}), collections.Permission([]models.Permission{{ID: 1}}))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(false))
		})
	})
})
