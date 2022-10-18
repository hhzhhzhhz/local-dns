package utils

import (
	"encoding/json"
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

func HttpSuccess(w http.ResponseWriter, data interface{}) error {
	b, _ := json.Marshal(&Response{
		Code: success,
		Data:  data,
	})
	w.Header().Add("content-type","application/json")
	_, err := w.Write(b)
	return err
}

func HttpFailed(w http.ResponseWriter, data constant.ApiResponse) error {
	b, _ := json.Marshal(&Response{
		Code: data.Code,
		Info:  data.Info,
	})
	w.Header().Add("content-type","application/json")
	_, err := w.Write(b)
	return err
}
