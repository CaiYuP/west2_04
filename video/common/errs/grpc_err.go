package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcError(err *BError) error {
	return status.Error(codes.Code(err.Code), err.Msg)
}
func ParseGrpcError(err error) (int32, string) {
	fromError, _ := status.FromError(err)
	return int32(fromError.Code()), fromError.Message()
}
func ToBError(err error) *BError {
	fromError, _ := status.FromError(err)
	return NewError(int32((fromError.Code())), fromError.Message())
}
