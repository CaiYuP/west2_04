package v1

import "time"

func Success() *BaseResponse {
	r := &BaseResponse{
		Code: SuccessCode,
		Msg:  "success",
	}
	return r
}
func (r *BaseResponse) Failed(code int32, msg string) {
	r.Code = code
	r.Msg = msg
}
func GetTimestamp(t time.Time) *Timestamp {
	return &Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}
