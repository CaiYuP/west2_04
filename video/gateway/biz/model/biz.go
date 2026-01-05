package model

var JwtSecret = []byte("your-secret-key")

const (
	AccountExist = iota
	AccountForbit
)

var RedisLogin = "user::"
var GetAll = 2
var (
	Failed = int32(10001)
)
