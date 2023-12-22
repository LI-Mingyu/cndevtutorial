#!/bin/bash

# 设置目标服务的URL
SERVICE_URL="http://localhost:8081"

# 设置请求的次数
REQUEST_COUNT=100

# 循环发送请求
for ((i=1; i<=REQUEST_COUNT; i++))
do
   curl $SERVICE_URL
   sleep 1
done

# 输出结束信息
echo "Completed $REQUEST_COUNT requests"
