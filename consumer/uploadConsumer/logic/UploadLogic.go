package logic

import (
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/gloabalType"
	"doushen_by_liujun/internal/util"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

func UploadHandle(_, message string) error {

	var uploadMessage gloabalType.UploadSuccessMessage
	if err := json.Unmarshal([]byte(message), &uploadMessage); err != nil {
		// 处理反序列化错误
		// ...
		logx.Error("消息传递类型错误")
	}
	// 持久log
	err := util.LogWrite(common.UploadSecurity, uploadMessage)
	if err != nil {
		return err
	}
	return nil
}
