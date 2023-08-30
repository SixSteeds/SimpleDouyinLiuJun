# SimpleDouyinLiuJun

一个基于go-zero的微服务简化版抖音项目，由六骏小队于23年8月参加字节跳动组织的第6届青训营时所做。

## 项目启动

#### 安装依赖

```
go mod tidy
```

#### 进入`script`目录安装本地镜像

```
docker compose -f docker-compose-env up -d
```
#### 服务启动
##### 进入`script`目录执行以下脚本
######  mac or linux
```bash
sh start_debug_LinuxMac.sh
```
###### windows 按以下顺序启动
```
所有consumer-job-contentrpc-其他rpc-所有api
```

## 业务架构图
![img.png](desc/img.png)

## 预览


## 贡献者


