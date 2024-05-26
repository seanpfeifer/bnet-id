package bnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/seanpfeifer/bnet-id/cookie"
)

const (
	cookieBNetID    = "bnet_id"
	cookieBattleTag = "battle_tag"

	// Perform an HTTP get on this endpoint using a player's authenticated http.Client to get their info
	bnetUserInfoURL = "https://oauth.battle.net/userinfo"
)

var ErrNoLogin = errors.New("no login cookie")

func ClearPlayerCookies(w http.ResponseWriter) {
	allCookies := []string{cookieBNetID, cookieBattleTag}
	for _, name := range allCookies {
		cookie.ExpireCookie(w, name)
	}
}

// PlayerFromRequest returns a Player struct from the cookies in the request, or an error if any cookie is missing.
func PlayerFromRequest(r *http.Request) (*Player, error) {
	bnetID, err := r.Cookie(cookieBNetID)
	if err != nil {
		return nil, ErrNoLogin
	}

	battleTag, err := r.Cookie(cookieBattleTag)
	if err != nil {
		return nil, ErrNoLogin
	}

	id, err := strconv.ParseInt(bnetID.Value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse BNetID: %w", err)
	}

	return &Player{
		BNetID:    ID(id),
		BattleTag: battleTag.Value,
	}, nil
}

// ID is our Battle.net ID that uniquely represents a player
type ID int64

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

type Player struct {
	BNetID    ID     `json:"id"`
	BattleTag string `json:"battletag"`
}

func (p Player) SetCookies(w http.ResponseWriter) {
	cookie.SetSecureCookie(w, cookieBNetID, p.BNetID.String())
	cookie.SetSecureCookie(w, cookieBattleTag, p.BattleTag)
}

func GetPlayerInfo(c *http.Client) (*Player, error) {
	resp, err := c.Get(bnetUserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var player Player
	err = json.Unmarshal(body, &player)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal player info: %w", err)
	}
	return &player, nil
}
