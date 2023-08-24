package logic

import (
	"doushen_by_liujun/internal/common"
	"doushen_by_liujun/internal/gloabalType"
	"doushen_by_liujun/internal/util"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

func LoginLogHandle(t, message string) error {
	var loginMessage gloabalType.LoginSuccessMessage
	if err := json.Unmarshal([]byte(message), &loginMessage); err != nil {
		// 处理反序列化错误
		// ...
		logx.Error("消息传递类型错误")
	}

	// 持久log
	err := util.LogWrite(common.USER_SECURITY, loginMessage)
	if err != nil {
		return err
	}

	return nil
}
