package restapi

import "net/http"

type response struct {
	Code    int         `json:"code" enums: "200, 404, 500"`
	Message string      `json:"message" example: "error"`
	Data    interface{} `json:"data" example: "data"`
}

func NewResponse(code int, message string, data interface{}) *response {
	return &response{code, message, data}
}

func GetOkResp(data interface{}) *response {
	return NewResponse(http.StatusOK, "ok", data)
}

func GetServerErrResp(err error) *response {
	return GetErrResp(http.StatusInternalServerError, err)
}

func GetErrResp(code int, err error) *response {
	return NewResponse(code, err.Error(), nil)
}
