package controller

import (
	"IM/common/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseData struct {
	Code    web.HttpStatusCode `json:"code"`
	Message interface{}        `json:"message"`
	Data    interface{}        `json:"data,omitempty"` // omitempty的作用是，当Data是空的时候，就不会展示出来了，有才会展示出来
}

func ResponseError(c *gin.Context, code web.HttpStatusCode, err interface{}) {
	rd := &ResponseData{
		Code:    code,
		Message: code.GetMsg(),
		Data:    err,
	}
	c.JSON(http.StatusOK, rd)
}

func ResponseSuccess(c *gin.Context, code web.HttpStatusCode, data interface{}) {
	rd := &ResponseData{
		Code:    code,
		Message: code.GetMsg(),
		Data:    data,
	}
	c.JSON(http.StatusOK, rd)
}
