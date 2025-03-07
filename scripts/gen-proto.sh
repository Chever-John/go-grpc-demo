#!/bin/bash

# 获取项目根目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." &> /dev/null && pwd )"

# 切换到项目根目录
cd "${PROJECT_ROOT}"

# 确保输出目录存在
mkdir -p pkg/pb/v1

echo "Generating protocol buffers code..."

# 生成 pb 文件
protoc --proto_path=. \
       --go_out=module=github.com/Chever-John/go-grpc-demo:. \
       --go-grpc_out=module=github.com/Chever-John/go-grpc-demo:. \
       api/proto/v1/*.proto

# 检查命令是否成功执行
if [ $? -eq 0 ]; then
    echo "Protocol buffers code generated successfully."
else
    echo "Error: Failed to generate protocol buffers code."
    exit 1
fi