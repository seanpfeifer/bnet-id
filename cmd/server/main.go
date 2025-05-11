package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/seanpfeifer/bnet-id/bnet"
	"github.com/seanpfeifer/rigging/fileload"
	"github.com/seanpfeifer/rigging/logging"
)

const (
	defaultSecretsFile = "secrets/secret.toml"
	host               = ":8080"
)

type Secrets struct {
	ClientID     string `toml:"clientID"`
	ClientSecret string `toml:"clientSecret"`
	RedirectURL  string `toml:"redirectURL"`
}

func main() {
	// Allow the user to specify a different secrets file
	secretsFile := flag.String("secrets", defaultSecretsFile, "Full path to the secrets file")
	flag.Parse()

	secrets, _, err := fileload.TOML[Secrets](*secretsFile)
	logging.FatalIfError(err, *secretsFile)

	srv := bnet.NewServer(secrets.ClientID, secrets.ClientSecret, secrets.RedirectURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.Index)
	mux.HandleFunc("/login", srv.OAuthRedirect)
	mux.HandleFunc("/oauthCallback", srv.OAuthCallback)
	mux.HandleFunc("/logout", srv.Logout)
	// A health check endpoint, just returns 200 OK if the web server is responding
	mux.HandleFunc("/health", func(http.ResponseWriter, *http.Request) {})
	// Static files, from the "static/" dir
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	httpSrv := &http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	httpSrv.ListenAndServe()
}
