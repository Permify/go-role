package repositories

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	`regexp`

	`github.com/Permify/permify-gorm/models`
	`github.com/Permify/permify-gorm/models/pivot`
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
				UserID:        1,
				RoleID:        1,
			}

			const sqlSelectOne = `SELECT count(*) FROM "user_roles" WHERE user_roles.user_id = $1 AND user_roles.role_id = $2`

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(userRoles.UserID, userRoles.RoleID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))

			db, err := repository.HasRole(1, models.Role{ID: 1})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(true))
		})
	})

})
