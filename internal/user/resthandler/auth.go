package resthandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/koteyye/news-portal/pkg/signer"
)

var errNoCookie = errors.New("отсутствует cookie авторизации, необходимо авторизоваться")

type ctxProfileKey string

const profileIDKey ctxProfileKey = "profile"

func (h RESTHandler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authorization")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				h.mapErrToResponse(w, http.StatusUnauthorized, errNoCookie)
				return
			}
			h.mapErrToResponse(w, http.StatusInternalServerError, fmt.Errorf("ошибка при получении cookie"))
			h.logger.Error(err.Error())
			return
		}
		profile, err := h.signer.Parse(cookie.Value)
		if err != nil {
			if errors.Is(err, signer.ErrTokenExpired) {
				h.mapErrToResponse(w, http.StatusUnauthorized, err)
				return
			}
			h.mapErrToResponse(w, http.StatusInternalServerError, fmt.Errorf("ошибка при парсинге токена"))
			h.logger.Error(err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), profileIDKey, profile)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
