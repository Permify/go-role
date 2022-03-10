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
})
