package repositories_scopes

import (
	"gorm.io/gorm"

	`github.com/Permify/permify-gorm/helpers`
	`github.com/Permify/permify-gorm/utils`
)

type GormPager interface {
	ToPaginate() func(db *gorm.DB) *gorm.DB
}

type GormPagination struct {
	*utils.Pagination
}

func (r *GormPagination) ToPaginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(helpers.OffsetCal(r.Pagination.GetPage(), r.Pagination.GetLimit())).Limit(r.Pagination.GetLimit())
	}
}
