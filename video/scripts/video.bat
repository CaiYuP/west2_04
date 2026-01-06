@echo off
REM =========================================
REM video.bat - 仅生成 video.proto 的 Go 代码
REM 1. 先切到项目根（脚本上级目录）
REM 2. 使用 paths=source_relative，保证生成文件与 .proto 同级
REM =========================================
echo 开始生成 video.proto ...

REM pushd 把当前目录压栈，然后进入脚本所在目录的上一级（项目根）
pushd "%~dp0.."

REM 运行 protoc
protoc ^
    --proto_path=. ^
    --go_out=paths=source_relative:. ^
    --go-grpc_out=paths=source_relative:. ^
    api/video/video.proto

if %errorlevel% neq 0 (
    echo 1
) else (
    echo 2
)

REM 回到之前目录
popd