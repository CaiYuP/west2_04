@echo off
REM Windows 批处理脚本：生成所有 gRPC 代码

echo 开始生成 gRPC 代码...

REM 生成 common
protoc --go_out=. --go-grpc_out=. --proto_path=. api/common/base.proto

REM 生成 user
protoc --go_out=. --go-grpc_out=. --proto_path=. api/user/user.proto

REM 生成 video
protoc --go_out=. --go-grpc_out=. --proto_path=. api/video/video.proto

REM 生成 interaction
protoc --go_out=. --go-grpc_out=. --proto_path=. api/interaction/interaction.proto

REM 生成 social
protoc --go_out=. --go-grpc_out=. --proto_path=. api/social/social.proto

echo gRPC 代码生成完成！





