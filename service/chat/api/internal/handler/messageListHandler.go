package handler

import (
	"fmt"
	"net/http"

	"doushen_by_liujun/service/chat/api/internal/logic"
	"doushen_by_liujun/service/chat/api/internal/svc"
	"doushen_by_liujun/service/chat/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func messageListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MessageChatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewMessageListLogic(r.Context(), svcCtx)
		resp, err := l.MessageList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			fmt.Println("in ok")
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
