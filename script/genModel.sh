#!/usr/bin/env bash

# 使用方法：
# ./model.sh usercenter user
# ./model.sh usercenter user_auth
# 再将./genModel下的文件剪切到对应服务的model目录里面，记得改package


#生成的表名
tables=follows,follows
#表生成的genmodel目录
modeldir=./model

# 数据库配置
host=127.0.0.1
port=3306
dbname=liujun_user
username=root
passwd=030321


goctl model mysql datasource -url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" -table="${tables}"  -dir="${modeldir}" -cache=true --style=goZero
goctl model mysql datasource -url="root:liujun@tcp(8.137.50.160:3306)/liujun_chat" -table="chat_message" -dir="./model" -cache=true --style=goZero
