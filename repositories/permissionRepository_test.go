package repositories

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"

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

			db, err := repository.GetPermissionByID(permission.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(permission))
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

			db, err := repository.GetPermissionByGuardName(permission.GuardName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(permission))
		})

		It("not found", func() {
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := repository.GetPermissionByGuardName("create-contact-permission")
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
	})

})
