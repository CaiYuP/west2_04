@echo off
echo 开始生成 user.proto...
protoc --go_out=. --go-grpc_out=. --proto_path=.     --go_out=paths=source_relative:.  --go-grpc_out=paths=source_relative:.  api/user/user.proto
if %errorlevel% neq 0 (
    echo 1
) else (
    echo 2
)