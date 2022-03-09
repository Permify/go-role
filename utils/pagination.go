package utils

type IPagination interface {
	Get() *Pagination
	GetPage() int
	GetLimit() int
}

type Pagination struct {
	Page  int
	Limit int
}

func (p *Pagination) Get() *Pagination {
	return p
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	return p.Limit
}
