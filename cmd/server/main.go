package main

import (
	"net/http"
	"time"

	"github.com/seanpfeifer/bnet-id/bnet"
	"github.com/seanpfeifer/rigging/fileload"
	"github.com/seanpfeifer/rigging/logging"
)

const (
	secretFile = "secrets/secret.toml"
	host       = ":8080"
)

type Secrets struct {
	ClientID     string `toml:"clientID"`
	ClientSecret string `toml:"clientSecret"`
	RedirectURL  string `toml:"redirectURL"`
}

func main() {
	secrets, _, err := fileload.TOML[Secrets](secretFile)
	logging.FatalIfError(err, secretFile)

	srv := bnet.NewServer(secrets.ClientID, secrets.ClientSecret, secrets.RedirectURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.Index)
	mux.HandleFunc("/login", srv.OAuthRedirect)
	mux.HandleFunc("/oauthCallback", srv.OAuthCallback)
	mux.HandleFunc("/logout", srv.Logout)

	httpSrv := &http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	httpSrv.ListenAndServe()
}
