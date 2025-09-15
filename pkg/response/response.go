package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageResponse struct {
	Response
	PageInfo PageInfo `json:"page_info"`
}

type PageInfo struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalCount  int `json:"total_count"`
	TotalPage   int `json:"total_page"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: msg,
		Data:    data,
	})
}

func Error(c *gin.Context, httpCode int, errCode int, msg string) {
	c.JSON(httpCode, Response{
		Code:    errCode,
		Message: msg,
		Data:    nil,
	})
}

func PageSuccess(c *gin.Context, data interface{}, pageInfo PageInfo) {
	c.JSON(http.StatusOK, PageResponse{
		Response: Response{
			Code:    200,
			Message: "success",
			Data:    data,
		},
		PageInfo: pageInfo,
	})
}
