package repositories_scopes

import (
	"gorm.io/gorm"

	`github.com/Permify/permify-gorm/helpers`
	`github.com/Permify/permify-gorm/utils`
)

// GormPager
// adds pagination capability to your gorm queries.
type GormPager interface {
	ToPaginate() func(db *gorm.DB) *gorm.DB
}

// GormPagination
// represent pagination data for pagination.
type GormPagination struct {
	*utils.Pagination
}

// ToPaginate
// adds pagination query to your gorm queries.
func (r *GormPagination) ToPaginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(helpers.OffsetCal(r.Pagination.GetPage(), r.Pagination.GetLimit())).Limit(r.Pagination.GetLimit())
	}
}
