package util

import (
	"bytes"
	"context"
	"doushen_by_liujun/internal/common"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"os"
)

// Upload 视频上传
func Upload(ctx context.Context, video []byte, fileName string) (*storage.PutRet, error) {

	fmt.Println("进入上传upload逻辑")
	putPolicy := storage.PutPolicy{
		Scope: common.BucketName,
	}
	mac := qbox.NewMac(common.AccessKey, common.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		UseHTTPS:      true,
		UseCdnDomains: true,
	}
	dataLen := int64(len(video))
	ret := storage.PutRet{}
	if dataLen > 5*1024*1024 {
		//采用分片上传（同时确保断点续传）
		resumeUploader := storage.NewResumeUploaderV2(&cfg)
		recorder, err := storage.NewFileRecorder(os.TempDir())
		if err != nil {
			return nil, err
		}
		putExtra := storage.RputV2Extra{
			Recorder: recorder,
			CustomVars: map[string]string{
				"x:name": "分片上传文件",
			},
		}
		reader := bytes.NewReader(video)
		err = resumeUploader.Put(ctx, &ret, upToken, fileName, reader, dataLen, &putExtra)
		if err != nil {
			return nil, err
		}
	} else {
		// 采用直接上传
		formUploader := storage.NewFormUploader(&cfg)
		ret := storage.PutRet{}
		putExtra := storage.PutExtra{
			Params: map[string]string{
				"x:name": "直传文件",
			},
		}
		err := formUploader.Put(context.Background(), &ret, upToken, fileName, bytes.NewReader(video), dataLen, &putExtra)
		if err != nil {
			return nil, err
		}
	}
	return &ret, nil
}
