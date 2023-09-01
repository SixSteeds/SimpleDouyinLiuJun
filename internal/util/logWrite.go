package util

import (
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strconv"
	"time"
)

func LogWrite(fileName string, T interface{}) error {

	jsonData, err := json.Marshal(T)
	if err != nil {
		// 处理序列化错误
		// ...
		logx.Error("序列化失败")
		return err
	}
	now := time.Now()
	year, month, day := now.Date()
	logDir := "../../../git logs/kafka-log/" + fileName + "/"
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		// 处理创建目录失败的错误
		fmt.Println("无法创建日志目录:", err)
		return err
	}
	filename := logDir + strconv.FormatInt(int64(year), 10) + "-" + strconv.FormatInt(int64(month), 10) + "-" + strconv.FormatInt(int64(day), 10) + ".json"

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// 处理打开文件错误
		return err
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
