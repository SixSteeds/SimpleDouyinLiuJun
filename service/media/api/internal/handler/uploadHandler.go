package handler

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"doushen_by_liujun/service/media/api/internal/logic"
	"doushen_by_liujun/service/media/api/internal/svc"
	"doushen_by_liujun/service/media/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func uploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req types.UploadReq
		_ = httpx.Parse(r, &req)

		file, _, err := r.FormFile("data")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("close file err: %v\n", err)
			}
		}(file)
		// 读取文件内容
		req.Data = make([]byte, 0)
		buf := make([]byte, 1024)
		for {
			n, err := file.Read(buf)
			if err != nil {
				break
			}
			req.Data = append(req.Data, buf[:n]...)
		}
		//fmt.Println(req)
		fmt.Println("uploadHandler err5555")
		ctx := context.WithValue(r.Context(), "ip", r.Header.Get("X-Real-IP"))

		l := logic.NewUploadLogic(ctx, svcCtx)
		resp, err := l.Upload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
