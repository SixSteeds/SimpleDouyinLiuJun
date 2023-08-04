package handler

import (
	"net/http"

	"doushen_by_liujun/service/chat/api/internal/logic"
	"doushen_by_liujun/service/chat/api/internal/svc"
	"doushen_by_liujun/service/chat/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func messageActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MessageActionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewMessageActionLogic(r.Context(), svcCtx)
		resp, err := l.MessageAction(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
