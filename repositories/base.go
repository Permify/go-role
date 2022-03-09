package repositories

// Seed

type Seedable interface {
	Seed() error
}

func Seeds(repos ...Seedable) (err error)  {
	for _, r := range repos {
		err = r.Seed()
	}
	return
}

// Migrate

type Migratable interface {
	Migrate() error
}

func Migrates(repos ...Migratable) (err error)  {
	for _, r := range repos {
		err = r.Migrate()
	}
	return
}

