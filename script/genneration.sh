

# api生成（更改.api后）
goctl api go -api *.api -dir ../  --style=goZero


# rpc生成（更改.proto后）
goctl rpc protoc *.proto --go_out=../ --go-grpc_out=../  --zrpc_out=../ --style=goZero



# model生成（更改数据库文件后）
sql2pb -go_package ./pb -host 43.136.243.177 -package pb -password liujun -port 3306 -schema liujun_user -service_name user -user root > user.proto



#将 `<topic-name>` 替换为您想要创建的主题名称， `<num-partitions>` 替换为主题的分区数，
kafka-topics.sh --bootstrap-server localhost:9092 --create --topic <topic-name>
# 创建topic
kafka-topics.sh --bootstrap-server localhost:9092 --create --topic loginLog
# 查看topic列表
kafka-topics.sh --bootstrap-server localhost:9092 --list


# 生成model
goctl model mysql datasource -url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" -table="${tables}"  -dir="${modeldir}" -cache=true --style=goZero
goctl model mysql datasource -url="root:liujun@tcp(8.137.50.160:3306)/liujun_chat" -table="chat_message" -dir="./model" -cache=true --style=goZero


# 例如 goctl docker --go user.go --exe userapi --version 1.19
goctl docker --go <service_name>.go --exe <docker_name> --version 1.19

# 例如 docker build -t userapi:v1 .
docker build -t <docker_name>:v1 .