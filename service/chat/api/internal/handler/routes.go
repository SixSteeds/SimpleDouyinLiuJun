// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"doushen_by_liujun/service/chat/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.JwtAuthMiddleware},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/action",
					Handler: messageActionHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/chat",
					Handler: messageListHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/douyin/message"),
	)
}
