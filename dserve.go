package dserve

import (
	"fmt"
	"net/http"
	"time"
)

// Serve launches HTTP server serving on listenAddr and servers a basic_auth secured directory at secure/static
func Serve(listenAddr string, secureDir bool, timeout time.Duration) error {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("."))
	mux.Handle("/", fs)

	if secureDir {
		if err := authInit(); err != nil {
			return fmt.Errorf("failed to initialize secure/ dir: %v", err)
		}
		mux.HandleFunc("/secure/", handleSecure)
	}

	svr := &http.Server{
		Addr:           listenAddr,
		Handler:        mux,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout * 2,
		MaxHeaderBytes: 1 << 20,
	}
	return svr.ListenAndServe()
}

func handleSecure(w http.ResponseWriter, r *http.Request) {
	if validBasicAuth(r) {
		fs := http.FileServer(http.Dir("secure/static"))
		h := http.StripPrefix("/secure/", fs)
		h.ServeHTTP(w, r)
		return
	}
	w.Header().Set("WWW-Authenticate", `Basic realm="Dserve secure/ Basic Authentication"`)
	http.Error(w, "Not Authorized", http.StatusUnauthorized)
}
