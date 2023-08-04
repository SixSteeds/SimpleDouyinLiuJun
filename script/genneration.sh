# goctl api go -api *.api -dir ../  --style=goZero
# goctl rpc protoc *.proto --go_out=../ --go-grpc_out=../  --zrpc_out=../ --style=goZero

# sql2pb -go_package ./pb -host 43.136.243.177 -package pb -password liujun -port 3306 -schema liujun_user -service_name user -user root > user.proto