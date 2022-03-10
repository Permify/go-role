package utils

// IPagination abstraction for your pagination data getters.
type IPagination interface {
	Get() *Pagination
	GetPage() int
	GetLimit() int
}

// Pagination pagination data.
type Pagination struct {
	Page  int
	Limit int
}

// Get get pagination struct.
func (p *Pagination) Get() *Pagination {
	return p
}

// GetPage get page value from pagination struct.
func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

// GetLimit get limit value from pagination struct.
func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	return p.Limit
}
