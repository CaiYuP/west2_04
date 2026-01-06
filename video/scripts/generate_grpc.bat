@echo off
REM Windows 批处理脚本：生成所有 gRPC 代码，保持生成文件与 proto 同级 (source_relative)

echo 开始生成 gRPC 代码...

REM 在项目根目录执行，输出到与 proto 同目录，避免生成额外的目录层级
protoc ^
    --go_out=paths=source_relative:. ^
    --go-grpc_out=paths=source_relative:. ^
    --proto_path=. ^
    api/common/base.proto ^

    api/video/video.proto ^
 api/user/user.proto ^
  api/interaction/interaction.proto ^
     api/social/social.proto

echo gRPC 代码生成完成！
