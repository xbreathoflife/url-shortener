package auth

import (
	"context"
	"github.com/xbreathoflife/url-shortener/internal/app/core"
	"net/http"
)

const cookieName = "uuid"
const CtxKey = ContextKey("uuid")

type ContextKey string


func AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		var uuid string
		if err != nil {
			uuid = core.GenerateUUID()
			encryptedUUID, err := core.Encrypt(uuid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			newCookie := http.Cookie{Name: "uuid", Value: encryptedUUID}
			http.SetCookie(w, &newCookie)
		} else {
			uuid, err  = core.Decrypt(cookie.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		ctx := context.WithValue(r.Context(), CtxKey, uuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
