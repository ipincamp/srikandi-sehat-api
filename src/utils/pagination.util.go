package utils

import (
	"ipincamp/srikandi-sehat/src/dto"
	"math"

	"gorm.io/gorm"
)

func GeneratePagination(page, limit int, db *gorm.DB, model interface{}) (dto.Pagination, func(db *gorm.DB) *gorm.DB) {
	var totalRows int64
	db.Model(model).Count(&totalRows)

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	var prevPage *int
	if page > 1 && page <= totalPages {
		p := page - 1
		prevPage = &p
	}

	var nextPage *int
	if page < totalPages {
		n := page + 1
		nextPage = &n
	}

	pagination := dto.Pagination{
		Limit:        limit,
		TotalRows:    totalRows,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PreviousPage: prevPage,
		NextPage:     nextPage,
	}

	paginateScope := func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}

	return pagination, paginateScope
}
