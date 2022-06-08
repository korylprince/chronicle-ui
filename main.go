package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
	auth "github.com/korylprince/go-ad-auth/v3"
)

func writeBadAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic")
	w.WriteHeader(http.StatusUnauthorized)
}

func WithAuth(config *Config, next http.Handler) http.Handler {
	authConfig := config.AuthConfig()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			writeBadAuth(w)
			return
		}

		stat, _, grps, err := auth.AuthenticateExtended(authConfig, u, p, nil, []string{config.LDAPGroup})
		if err != nil {
			log.Println("could not authenticate:", err)
			writeBadAuth(w)
			return
		}

		if !stat || len(grps) == 0 {
			writeBadAuth(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RunServer() error {
	config := new(Config)
	err := envconfig.Process("", config)
	if err != nil {
		return fmt.Errorf("could not process configuration from environment: %w", err)
	}

	db, err := NewDB(config.SQLDSN)
	if err != nil {
		return fmt.Errorf("could not create db: %w", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		form, err := db.HandleUI(r)
		if err != nil {
			log.Println("Error:", err)
			form = &Form{Error: err}
		}

		w.WriteHeader(http.StatusOK)
		if err = Tmpl.Execute(w, form); err != nil {
			log.Println("could not execute template:", err)
		}
	})

	var handler http.Handler = handlers.CombinedLoggingHandler(os.Stdout, WithAuth(config, mux))
	// rewrite for x-forwarded-for, etc headers
	if config.ProxyHeaders {
		handler = handlers.ProxyHeaders(handler)
	}

	return http.ListenAndServe(config.ListenAddr, handler)
}

func main() {
	if err := RunServer(); err != nil {
		log.Println("could not start server:", err)
	}
}
