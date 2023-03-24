package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Result  bool        `json:"result"`
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponseError(result bool, ctx *gin.Context, status int, message string) {
	res := Response{
		Result:  result,
		Status:  status,
		Message: message,
	}
	ctx.JSON(status, res)
}

func NewResponseSuccessWithData(result bool, ctx *gin.Context, status int, message string, data interface{}) {
	res := Response{
		Result:  result,
		Status:  status,
		Message: message,
		Data:    data,
	}
	ctx.JSON(status, res)
}

func NewResponseSuccess(result bool, ctx *gin.Context, status int, message string) {
	res := Response{
		Result:  result,
		Status:  status,
		Message: message,
	}
	ctx.JSON(status, res)
}
