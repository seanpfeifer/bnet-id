package bnet

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"html/template"
	"time"

	"fmt"
	"net/http"

	"github.com/seanpfeifer/rigging/logging"
	"golang.org/x/oauth2"
)

const (
	stateValidTime   = 5 * time.Minute
	randomByteLength = 12 // BNet doesn't like large states, so I'm using 12 here
)

//go:embed site.tmpl
var siteTemplate string

var usOAuthEndpoint = oauth2.Endpoint{
	AuthURL:  "https://oauth.battle.net/authorize",
	TokenURL: "https://oauth.battle.net/token",
}

type Server struct {
	Nonces   NonceMap
	OAuthCfg *oauth2.Config
	site     *template.Template
}

func NewServer(clientID, clientSecret, redirectURL string) *Server {
	tmpl, err := template.New("site").Parse(siteTemplate)
	logging.FatalIfError(err)
	return &Server{
		Nonces: *NewNonceMap(stateValidTime),
		OAuthCfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     usOAuthEndpoint,
		},
		site: tmpl,
	}
}

func (s *Server) OAuthRedirect(w http.ResponseWriter, r *http.Request) {
	nonce := newNonce()
	s.Nonces.Add(nonce, stateValidTime)
	http.Redirect(w, r, s.OAuthCfg.AuthCodeURL(nonce), http.StatusFound)
}

func (s *Server) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	// Also check state
	state := r.URL.Query().Get("state")
	if !s.Nonces.Remove(state) {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	player, err := s.GetAccountInfo(code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get player info: %v", err), http.StatusInternalServerError)
		return
	}

	player.SetCookies(w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	player, _ := PlayerFromRequest(r)
	s.renderPage(w, player)
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	ClearPlayerCookies(w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) GetAccountInfo(code string) (*Player, error) {
	playerToken, err := s.OAuthCfg.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Create a client we can use with the player's token
	client := s.OAuthCfg.Client(context.Background(), playerToken)

	// Finally, we can retrieve the player's info
	return GetPlayerInfo(client)
}

func (s *Server) renderPage(w http.ResponseWriter, player *Player) {
	if err := s.site.Execute(w, player); err != nil {
		http.Error(w, fmt.Sprintf("Failed to render page: %v", err), http.StatusInternalServerError)
	}
}

func newNonce() string {
	var b [randomByteLength]byte
	_, err := rand.Read(b[:])
	logging.FatalIfError(err)

	return base64.RawURLEncoding.EncodeToString(b[:])
}
