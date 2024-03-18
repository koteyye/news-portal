package resthandler

import (
	"errors"
	"github.com/koteyye/news-portal/pkg/models"
	resp "github.com/koteyye/news-portal/pkg/restresponser"
	"net/http"
)

const writerRole = "writer"

var (
	errNoProfile = errors.New("profile is empty")
	errNoAllowed = errors.New("not allowed")
)

func (h RESTHandler) checkWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile := r.Context().Value(profileIDKey).(*models.Profile)
		if profile == nil {
			resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusForbidden, Err: errNoProfile, ContentType: resp.CtJSON})
			return
		}
		allowed := false
		for _, role := range profile.Roles {
			if role == writerRole {
				allowed = true
			}
		}
		if !allowed {
			resp.MapErrToResponse(w, &resp.ResponseOptions{StatusCode: http.StatusForbidden, Err: errNoAllowed, ContentType: resp.CtJSON})
			return
		}
		next.ServeHTTP(w, r)
	})
}
