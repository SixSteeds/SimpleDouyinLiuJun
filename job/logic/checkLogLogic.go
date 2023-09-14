package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"path/filepath"
	"time"
)

type CheckLogLogic struct {
	ctx context.Context
	logx.Logger
}

func NewCheckLogLogic(ctx context.Context) *CheckLogLogic {
	return &CheckLogLogic{
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckLogLogic) CheckLog() {
	/*
		Author：    刘洋
		Function：  每天点定时扫描 log/uploadSecurity 目录和 log/userSecurity 目录，超过7天的文件清除
		Update：    08.25
	*/

	var dirs = make([]string, 0)
	dirPath, _ := os.Getwd()
	dirs = append(dirs, filepath.Join(dirPath, "/logs/kafka-log/uploadSecurity"))
	dirs = append(dirs, filepath.Join(dirPath, "/logs/kafka-log/userSecurity"))

	// 遍历所有要检查log的目录
	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fileName := info.Name()
				fileTime := fileName[:len(fileName)-5]
				//fmt.Println(fileTime)
				fileDate, err := time.Parse("2006-1-2", fileTime)
				if err != nil {
					fmt.Println("Failed to parse string:", err)
					return nil
				}
				// 如果是文件且文件名符合日期格式，则比较日期并删除7天前的文件
				if time.Now().Sub(fileDate) > 7*24*time.Hour {
					fmt.Printf("Deleting file: %s\n", path)
					err := os.Remove(path)
					if err != nil {
						fmt.Println("Failed to delete file:", err)
					}
				}
			}
			fmt.Println("do nothing")
			return nil
		})
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
	fmt.Println("【checkLog并更新完成（7天以上记录清除）】")
	return
}
