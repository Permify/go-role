package helpers

import (
	`math`
)

func NextPageCal(page int, totalPage int) int {
	if page == totalPage {
		return page
	}
	return page + 1
}

func PrevPageCal(page int) int {
	if page > 1 {
		return page - 1
	}
	return page
}

func TotalPage(count int64, limit int) int {
	return int(math.Ceil(float64(count) / float64(limit)))
}

func OffsetCal(page int, limit int) int {
	return (page - 1) * limit
}
