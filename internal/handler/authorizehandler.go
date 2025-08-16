package handler

import (
	"net/http"

	"oauth2-server/internal/logic"
	"oauth2-server/internal/svc"
	"oauth2-server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AuthorizeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthorizeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAuthorizeLogic(r.Context(), svcCtx)
		resp, err := l.Authorize(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
