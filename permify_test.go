package permify_gorm

import (
	`database/sql`
	`github.com/DATA-DOG/go-sqlmock`
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	`gorm.io/driver/postgres`
	`gorm.io/gorm`
	`regexp`
	"testing"

	`github.com/Permify/permify-gorm/helpers`
	`github.com/Permify/permify-gorm/models`
	`github.com/Permify/permify-gorm/repositories`
)

var permify *Permify

func TestPermify(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Permify suite")
}

var _ = Describe("Permify Service", func() {
	var permify *Permify
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


		roleRepository := &repositories.RoleRepository{Database: gormDb}
		permissionRepository := &repositories.PermissionRepository{Database: gormDb}
		userRepository := &repositories.UserRepository{Database: gormDb}

		permify = &Permify{
			RoleRepository:       roleRepository,
			PermissionRepository: permissionRepository,
			UserRepository:       userRepository,
		}
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

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE permissions.id = $1 ORDER BY "permissions"."id" LIMIT 1`)).
				WithArgs(permission.ID).
				WillReturnRows(rows)

			db, err := permify.GetPermission(permission.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(permission))
		})

		It("not found", func() {
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := permify.GetPermission(1)
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
	})

	Context("Get Permission By Name", func() {
		It("found", func() {
			permission := models.Permission{
				ID:        1,
				Name:      "create contact permission",
				GuardName: "create-contact-permission",
			}

			rows := sqlmock.NewRows([]string{"id", "name", "guard_name"}).
				AddRow(permission.ID, permission.Name, permission.GuardName)

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE permissions.guard_name = $1 ORDER BY "permissions"."id" LIMIT 1`)).
				WithArgs(helpers.Guard(permission.Name)).
				WillReturnRows(rows)

			db, err := permify.GetPermission(helpers.Guard(permission.Name))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(permission))
		})

		It("not found", func() {
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := permify.GetPermission("create-contact-permission")
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
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

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE roles.id = $1 ORDER BY "roles"."id" LIMIT 1`)).
				WithArgs(role.ID).
				WillReturnRows(rows)

			db, err := permify.GetRole(role.ID, false)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(role))
		})

		It("not found", func() {
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := permify.GetRole(1, false)
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
	})

	Context("Get Permission By Name", func() {
		It("found", func() {
			role := models.Role{
				ID:        1,
				Name:      "admin",
				GuardName: "admin",
			}

			rows := sqlmock.NewRows([]string{"id", "name", "guard_name"}).
				AddRow(role.ID, role.Name, role.GuardName)

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE roles.guard_name = $1 ORDER BY "roles"."id" LIMIT 1`)).
				WithArgs(helpers.Guard(role.Name)).
				WillReturnRows(rows)

			db, err := permify.GetRole(helpers.Guard(role.Name), false)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(db).Should(Equal(role))
		})

		It("not found", func() {
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := permify.GetRole("admin", false)
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
	})
})
