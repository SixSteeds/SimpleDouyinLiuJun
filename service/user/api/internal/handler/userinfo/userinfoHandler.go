package userinfo

import (
	"net/http"

	"doushen_by_liujun/service/user/api/internal/logic/userinfo"
	"doushen_by_liujun/service/user/api/internal/svc"
	"doushen_by_liujun/service/user/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserinfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserinfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := userinfo.NewUserinfoLogic(r.Context(), svcCtx)
		resp, err := l.Userinfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
