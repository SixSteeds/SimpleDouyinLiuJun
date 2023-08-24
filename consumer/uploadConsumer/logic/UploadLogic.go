package logic

import (
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/gloabalType"
	"doushen_by_liujun/internal/util"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func UploadHandle(t, message string) error {
	var uploadMessage gloabalType.UploadSuccessMessage
	if err := json.Unmarshal([]byte(message), &uploadMessage); err != nil {
		// 处理反序列化错误
		// ...
		logx.Error("消息传递类型错误")
	}

	// 1、拿到封面

	// 2、写入数据库
	conn := sqlx.NewMysql("root:liujun@tcp(8.137.50.160:3306)/liujun_content?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
	_ = conn

	// 持久log
	err := util.LogWrite(common.UPLOAD_SECURITY, uploadMessage)
	if err != nil {
		return err
	}
	return nil
}
