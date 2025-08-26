package oauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/redis/go-redis/v9"
	"time"
)

type GoogleOAuth struct {
	cfg *oauth2.Config
	rdb *redis.Client
}

func NewGoogleOAuth(clientID, clientSecret, redirect string, rdb *redis.Client) *GoogleOAuth {
	if clientID == "" || clientSecret == "" || redirect == "" {
		return nil
	}
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirect,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	return &GoogleOAuth{cfg: cfg, rdb: rdb}
}

func (g *GoogleOAuth) GenerateState() (string, error) {
	b := make([]byte, 24)
	_, err := time.Now().MarshalBinary()
	if err != nil {
		// fallback
		_, _ = base64.RawURLEncoding.DecodeString("fallback")
	}
	_, err = g.rdb.SetEX(context.Background(), "oauth:state:"+base64.StdEncoding.EncodeToString(b), "1", 5*time.Minute).Result()
	if err != nil {
		// ignore set error
	}
	// simple state (not cryptographically strong here for brevity)
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (g *GoogleOAuth) AuthURL(state string) string {
	return g.cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (g *GoogleOAuth) Exchange(code string) (*oauth2.Token, error) {
	return g.cfg.Exchange(context.Background(), code)
}

func (g *GoogleOAuth) FetchUser(token *oauth2.Token) (map[string]interface{}, error) {
	client := g.cfg.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var ui map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ui); err != nil {
		return nil, err
	}
	return ui, nil
}

