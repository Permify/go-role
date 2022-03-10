package repositories

// Seedable gives models the ability to seed.
type Seedable interface {
	Seed() error
}

// Seeds seed to seedable models.
func Seeds(repos ...Seedable) (err error) {
	for _, r := range repos {
		err = r.Seed()
	}
	return
}

// Migratable gives models the ability to migrate.
type Migratable interface {
	Migrate() error
}

// Migrates migrate to migratable models.
func Migrates(repos ...Migratable) (err error) {
	for _, r := range repos {
		err = r.Migrate()
	}
	return
}
