package graph

import (
	"context"
	"net/http"
)

func InjectIDCardToCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		IDCard := r.Header.Get("Authorization")
		ctx := context.WithValue(r.Context(), "IDCard", IDCard)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
