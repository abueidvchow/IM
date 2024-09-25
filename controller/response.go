package controller

import (
	"IM/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseData struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data,omitempty"` // omitempty的作用是，当Data是空的时候，就不会展示出来了，有才会展示出来
}

func ResponseError(c *gin.Context, code common.HttpStatusCode, err ...interface{}) {
	rd := &ResponseData{
		Code:    code,
		Message: code.GetMsg(),
		Data:    err,
	}
	c.JSON(http.StatusOK, rd)
}

func ResponseSuccess(c *gin.Context, code common.HttpStatusCode, data ...interface{}) {
	rd := &ResponseData{
		Code:    code,
		Message: code.GetMsg(),
		Data:    data,
	}
	c.JSON(http.StatusOK, rd)
}
