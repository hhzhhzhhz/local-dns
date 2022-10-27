package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/hhzhhzhhz/local-dns/constant"
	"net/http"
)

const (
	success = 0
)

type Response struct {
	Code int `json:"code"`
	Info string `json:"info"`
	Data interface{} `json:"data"`
}

func HttpSuccess(c *gin.Context, data interface{}) {
	c.Header("content-type", "application/json")
	c.JSON(http.StatusOK, &Response{Code: success, Data:  data})
}

func HttpFailed(c *gin.Context, data constant.ApiResponse) {
	c.Header("content-type", "application/json")
	c.JSON(http.StatusOK, &Response{Code: data.Code, Info:  data.Info})
}
