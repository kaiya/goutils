package router

import (
	"net/http"
	"testing"

	"gitlab.momoso.com/cm/kit/third_party/lg"
)

func TestRoute(t *testing.T) {
	r := NewRouter()
	r.Use(logger)
	r.Use(ratelimit)
	r.Add("/hello", http.HandlerFunc(helloHandler))
	lg.PanicError(http.ListenAndServe(":8088", r))
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lg.Infof("before logger")
		next.ServeHTTP(w, r)
		lg.Infof("after logger")
	})
}

func ratelimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lg.Infof("before ratelimiter")
		next.ServeHTTP(w, r)
		lg.Infof("after ratelimiter")
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
	lg.Info("helloHandler finished")
}
