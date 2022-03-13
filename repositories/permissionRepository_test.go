package repositories

import (
	"database/sql"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Permify/permify-gorm/models"
)

var _ = Describe("Permission Repository", func() {
	var repository *PermissionRepository
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

		repository = &PermissionRepository{Database: gormDb}
	})

	AfterEach(func() {
		err := mock.ExpectationsWereMet()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("Get Permission By ID", func() {
		It("found", func() {
			permission := models.Permission{
				ID:        1,
				Name:      "create contact permission",
				GuardName: "create-contact-permission",
			}

			rows := sqlmock.NewRows([]string{"id", "name", "guard_name"}).
				AddRow(permission.ID, permission.Name, permission.GuardName)

			const sqlSelectOne = `SELECT * FROM "permissions" WHERE permissions.id = $1 ORDER BY "permissions"."id" LIMIT 1`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(permission.ID).
				WillReturnRows(rows)

			value, err := repository.GetPermissionByID(permission.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(permission))
		})

		It("not found", func() {
			// ignore sql match
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := repository.GetPermissionByID(1)
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
	})

	Context("Get Permission By Guard Name", func() {
		It("found", func() {
			permission := models.Permission{
				ID:        1,
				Name:      "create contact permission",
				GuardName: "create-contact-permission",
			}

			rows := sqlmock.NewRows([]string{"id", "name", "guard_name"}).
				AddRow(permission.ID, permission.Name, permission.GuardName)

			const sqlSelectOne = `SELECT * FROM "permissions" WHERE permissions.guard_name = $1 ORDER BY "permissions"."id" LIMIT 1`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(permission.GuardName).
				WillReturnRows(rows)

			value, err := repository.GetPermissionByGuardName(permission.GuardName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(permission))
		})

		It("not found", func() {
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := repository.GetPermissionByGuardName("create-contact-permission")
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
	})

	Context("Get Permissions", func() {
		It("found", func() {
			permissions := []models.Permission{
				{
					ID:        1,
					Name:      "create contact permission",
					GuardName: "create-contact-permission",
				},
			}

			rows := sqlmock.NewRows([]string{"id", "name", "guard_name"}).
				AddRow(permissions[0].ID, permissions[0].Name, permissions[0].GuardName)

			const sqlSelectOne = `SELECT * FROM "permissions" WHERE permissions.id IN ($1)`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(permissions[0].ID).
				WillReturnRows(rows)

			value, err := repository.GetPermissions([]uint{permissions[0].ID})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value.Origin()).Should(Equal(permissions))
		})
	})

	Context("Get Permissions By Guard Names", func() {
		It("found", func() {
			permissions := []models.Permission{
				{
					ID:        1,
					Name:      "create contact permission",
					GuardName: "create-contact-permission",
				},
			}

			rows := sqlmock.NewRows([]string{"id", "name", "guard_name"}).
				AddRow(permissions[0].ID, permissions[0].Name, permissions[0].GuardName)

			const sqlSelectOne = `SELECT * FROM "permissions" WHERE permissions.guard_name IN ($1)`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(permissions[0].GuardName).
				WillReturnRows(rows)

			value, err := repository.GetPermissionsByGuardNames([]string{permissions[0].GuardName})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value.Origin()).Should(Equal(permissions))
		})
	})
})
