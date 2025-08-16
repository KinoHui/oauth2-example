package handler

import (
	"net/http"

	"oauth2-server/internal/logic"
	"oauth2-server/internal/svc"
	"oauth2-server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ClientRegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ClientRegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewClientRegisterLogic(r.Context(), svcCtx)
		resp, err := l.ClientRegister(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
