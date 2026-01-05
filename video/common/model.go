package common

import "net/http"

type BusinessCode int
type Result struct {
	Code BusinessCode `json:"code" form:"code"`
	Msg  string       `json:"msg"  form:"msg"`
	Data interface{}  `json:"data"  form:"data"`
}

func (r *Result) Success(data interface{}) *Result {
	r.Data = data
	r.Code = http.StatusOK
	r.Msg = "success"
	return r
}
func (r *Result) Fail(code BusinessCode, msg string) *Result {
	r.Code = code
	r.Msg = msg
	return r
}
