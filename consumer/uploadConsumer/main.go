package main

import (
	"doushen_by_liujun/consumer/uploadConsumer/logic"
	"flag"
	"fmt"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "consumer/uploadConsumer/etc/uploadConsumer.yaml", "the config file")

type KqConf struct {
	service.ServiceConf
	Brokers     []string
	Group       string
	Topic       string
	Offset      string `json:",options=first|last,default=last"`
	Conns       int    `json:",default=1"`
	Consumers   int    `json:",default=8"`
	Processors  int    `json:",default=8"`
	MinBytes    int    `json:",default=10240"`    // 10K
	MaxBytes    int    `json:",default=10485760"` // 10M
	Username    string `json:",optional"`
	Password    string `json:",optional"`
	ForceCommit bool   `json:",default=true"`
}

func main() {
	var c KqConf
	conf.MustLoad(*configFile, &c)
	q := kq.MustNewQueue(kq.KqConf(c), kq.WithHandle(logic.UploadHandle))
	defer q.Stop()
	fmt.Println("Starting upload consumer...")
	q.Start()
}
