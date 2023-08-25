#!/bin/bash

# 定义要启动的Go文件路径列表
files=(
    "service/content/rpc/content.go"
    "service/user/rpc/user.go"
    "service/media/rpc/media.go"
    "service/chat/rpc/chat.go"

    "service/user/api/user.go"
    "service/content/api/content.go"
    "service/media/api/media.go"
    "service/chat/api/chat.go"

    "consumer/uploadConsumer/main.go"
    "consumer/loginConsumer/main.go"
    "job/main.go"
)

rm -rf debug_log
mkdir debug_log

# 循环启动Go文件并将输出写入独立的文件
for file in "${files[@]}"
do
    echo "启动服务: $file"
    output_file="debug_log/$(echo "$file" | awk -F/ '{print $(NF-1),$NF}' | cut -d. -f1 | awk '{print $1"_"$2}').log"
    # 创建output_file 文件
    echo "$output_file"
    touch $output_file
    nohup go run "$file" > "$output_file" 2>&1 &
    sleep 5
done

# 等待所有进程结束
wait
