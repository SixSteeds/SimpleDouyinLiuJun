
# api生成（更改.api后）
# goctl api go -api *.api -dir ../  --style=goZero

# rpc生成（更改.proto后）
# goctl rpc protoc *.proto --go_out=../ --go-grpc_out=../  --zrpc_out=../ --style=goZero


# model生成（更改数据库文件后）
# sql2pb -go_package ./pb -host 43.136.243.177 -package pb -password liujun -port 3306 -schema liujun_user -service_name user -user root > user.proto



#将 `<topic-name>` 替换为您想要创建的主题名称， `<num-partitions>` 替换为主题的分区数， `<replication-factor>` 替换为主题的副本数。
#kafka-topics.sh --bootstrap-server localhost:9092 --create --topic <topic-name> --partitions <num-partitions> --replication-factor <replication-factor>







