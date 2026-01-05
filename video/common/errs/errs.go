package errs

import (
	"fmt"
	"west2-video/common"
)

type BError struct {
	Code int32  `json:"code"`
	Msg  string `json:"message"`
}

func (e *BError) Error() string {
	return fmt.Sprintf("code:%d, msg:%s", e.Code, e.Msg)
}
func NewError(code int32, msg string) *BError {
	return &BError{
		Code: code,
		Msg:  msg,
	}
}
func (e *BError) ParseError() (common.BusinessCode, string) {
	return common.BusinessCode(e.Code), e.Msg
}
func ParseError(err error) (bool, common.BusinessCode, string) {
	if err != nil {
		switch err.(type) {
		case *BError:
			berr := err.(*BError)
			if berr != nil {
				parseError, s := berr.ParseError()
				return true, parseError, s
			}
		}
	}
	return false, common.BusinessCode(0), ""
}
