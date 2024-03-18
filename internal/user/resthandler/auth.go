package resthandler

import (
	"context"
	"errors"
	"net/http"

	resp "github.com/koteyye/news-portal/pkg/restresponser"
	"github.com/koteyye/news-portal/pkg/signer"
)

var errNoCookie = errors.New("auth cookie is empty")

type ctxProfileKey string

const profileIDKey ctxProfileKey = "profile"

func (h RESTHandler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authorization")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusUnauthorized, Err: errNoCookie, ContentType: resp.CtJSON})
				return
			}
			resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusInternalServerError, Err: errors.New("can't get cookie"), ContentType: resp.CtJSON})
			h.logger.Error(err.Error())
			return
		}
		profile, err := h.signer.Parse(cookie.Value)
		if err != nil {
			if errors.Is(err, signer.ErrTokenExpired) {
				resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusUnauthorized, Err: err, ContentType: resp.CtJSON})
				return
			}
			resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusInternalServerError, Err: errors.New("token parse error"), ContentType: resp.CtJSON})
			h.logger.Error(err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), profileIDKey, profile)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
