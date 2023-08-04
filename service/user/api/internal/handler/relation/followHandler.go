package relation

import (
	"doushen_by_liujun/service/user/api/internal/types"
	"net/http"

	"doushen_by_liujun/service/user/api/internal/logic/relation"
	"doushen_by_liujun/service/user/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FollowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		println("FollowHandler")
		var req types.FollowReq
		//if err := httpx.Parse(r, &req); err != nil {
		//	httpx.ErrorCtx(r.Context(), w, err)
		//	return
		//}

		l := relation.NewFollowLogic(r.Context(), svcCtx)
		resp, err := l.Follow(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
