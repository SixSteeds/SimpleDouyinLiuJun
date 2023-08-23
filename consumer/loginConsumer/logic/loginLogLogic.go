package logic

import (
	"doushen_by_liujun/internal/gloabalType"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

func LoginLogHandle(time, message string) error {
	var loginMessage gloabalType.LoginSuccessMessage
	fmt.Println(time)
	fmt.Println("===========================")
	fmt.Println(message)

	if err := json.Unmarshal([]byte(message), &loginMessage); err != nil {
		// 处理反序列化错误
		// ...
		logx.Error("消息传递类型错误")
	}
	fmt.Println(loginMessage.IP)
	fmt.Println(loginMessage.Logintime)
	fmt.Println(loginMessage.UserId)
	fmt.Println(loginMessage)
	return nil
}
