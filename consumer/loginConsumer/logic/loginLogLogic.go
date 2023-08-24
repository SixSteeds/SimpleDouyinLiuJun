package logic

import (
	"doushen_by_liujun/internal/gloabalType"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strconv"
	"time"
)

func LoginLogHandle(t, message string) error {
	var loginMessage gloabalType.LoginSuccessMessage

	if err := json.Unmarshal([]byte(message), &loginMessage); err != nil {
		// 处理反序列化错误
		// ...
		logx.Error("消息传递类型错误")
	}

	jsonData, err := json.Marshal(loginMessage)
	if err != nil {
		// 处理序列化错误
		// ...
		logx.Error("序列化失败")
		return err
	}
	now := time.Now()
	year, month, day := now.Date()

	filename := "consumer/loginConsumer/" + strconv.FormatInt(int64(year), 10) + "-" + strconv.FormatInt(int64(month), 10) + "-" + strconv.FormatInt(int64(day), 10) + ".json"

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// 处理打开文件错误
	}
	defer file.Close()

	if _, err := file.Write(jsonData); err != nil {
		logx.Error("写入文件失败")
		return err
	}
	// 在写入完毕后插入换行符
	if _, err := file.WriteString("\n"); err != nil {
		logx.Error("写入换行符失败")
		return err
	}

	return nil
}
