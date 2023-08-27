package middleware

import (
	"doushen_by_liujun/internal/util"
	"net/http"
)

type JwtAuthMiddleware struct {
}

func NewJwtAuthMiddleware() *JwtAuthMiddleware {
	return &JwtAuthMiddleware{}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		_, err := util.ParseToken(token)
		if err != nil {
			// 受到抖声app不携带token乱发请求的限制，注释此句
			//httpx.ErrorCtx(r.Context(), w, err)
		}
		//Pass through to next handler
		next(w, r)
	}
}
