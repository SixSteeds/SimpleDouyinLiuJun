package middleware

import (
	"doushen_by_liujun/internal/util"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type JwtAuthMiddleware struct {
}

func NewJwtAuthMiddleware() *JwtAuthMiddleware {
	return &JwtAuthMiddleware{}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		token := r.URL.Query().Get("token")
		mc, err := util.ParseToken(token)
		logx.Error(mc.Username)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		//Pass through to next handler
		next(w, r)
	}
}
