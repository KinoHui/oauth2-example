package handler

import (
	"net/http"

	"oauth2-server/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/client/register",
				Handler: ClientRegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/oauth/authorize",
				Handler: AuthorizeHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/oauth/token",
				Handler: TokenHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/oauth/userinfo",
				Handler: UserInfoHandler(serverCtx),
			},
		},
	)
}
