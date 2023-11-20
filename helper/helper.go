package helper

import "gorm.io/gorm"

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type Response struct {
	Meta      Meta        `json:"meta"`
	Data      interface{} `json:"data"`
	PerPage   int         `json:"per_page"`
	Page      int         `json:"page"`
	TotalData int64       `json:"total_data"`
	TotalPage int         `json:"total_page"`
}

func APIResponse(message string, code int, status string, data interface{}, perPage int, page int, totalData int64, totalPage int) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonResponse := Response{
		Meta:      meta,
		Data:      data,
		PerPage:   perPage,
		Page:      page,
		TotalData: totalData,
		TotalPage: totalPage,
	}
	return jsonResponse
}

func PaginationScopes(page int, perPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * perPage
		return db.Offset(offset).Limit(perPage)
	}
}
