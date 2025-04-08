#!/bin/bash

SERVER="localhost:50054"
USE_TLS=false
CA_FILE="ca/ca.crt"
SERVER_NAME="a.grpc.test.com"

if [ "$USE_TLS" = true ]; then
    TLS_OPT="-cacert $CA_FILE -servername $SERVER_NAME"
else
    TLS_OPT="-plaintext"
fi

echo "列出所有服务："
grpcurl $TLS_OPT $SERVER list

echo -e "\n描述 Demo 服务："
grpcurl $TLS_OPT $SERVER describe demo.Demo

echo -e "\n测试 Add 方法："
grpcurl $TLS_OPT -d '{"x": 10, "y": 2}' $SERVER demo.Demo/Add

echo -e "\n测试 SayHello 方法："
grpcurl $TLS_OPT -d '{"name": "张三"}' $SERVER demo.Demo/SayHello

echo -e "\n测试 GetStream 方法："
grpcurl $TLS_OPT -d '{"x": 10, "y": 2}' $SERVER demo.Demo/GetStream

echo -e "\n测试 PutStream 方法："
echo '{"x": 1}
{"x": 2}
{"x": 3}
{"x": 4}' | grpcurl $TLS_OPT -d @ $SERVER demo.Demo/PutStream

echo -e "\n测试 DoubleStream 方法："
echo '{"x": 0, "y": 0}
{"x": 1, "y": 1}
{"x": 2, "y": 2}
{"x": 3, "y": 3}' | grpcurl $TLS_OPT -d @ $SERVER demo.Demo/DoubleStream
