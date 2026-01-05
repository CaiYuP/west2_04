@echo off
echo 开始生成 user.proto...
protoc --go_out=. --go-grpc_out=. --proto_path=. api/user/user.proto
if %errorlevel% neq 0 (
    echo 生成失败！错误代码: %errorlevel%
    exit /b %errorlevel%
) else (
    echo 生成成功！文件位置: west2-video\api\user\v1\
)