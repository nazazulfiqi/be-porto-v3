package dto

import (
	"github.com/gin-gonic/gin"
)

// Struct untuk format response API umum
type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"` // Data bisa kosong (omitempty)
}

// Struct untuk response dengan pagination
type PaginationResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Struct untuk pagination
type Pagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	MaxPage   int `json:"max_page"`
	TotalData int `json:"total_data"`
}

// Helper untuk response sukses tanpa pagination
func SuccessResponse(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

// Helper untuk response error
func ErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       nil,
	})
}

// Helper untuk response sukses dengan pagination
func SuccessPaginationResponse(ctx *gin.Context, statusCode int, message string, data interface{}, page, limit, maxPage, totalData int) {
	ctx.JSON(statusCode, PaginationResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Pagination: Pagination{
			Page:      page,
			Limit:     limit,
			MaxPage:   maxPage,
			TotalData: totalData,
		},
	})
}
