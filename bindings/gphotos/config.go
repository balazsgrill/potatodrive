package gphotos

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/balazsgrill/potatodrive/gpfs"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	ClientID     string `json:"client_id" flag:"client-id" reg:"ClientID"`
	ClientSecret string `json:"client_secret" flag:"client-secret" reg:"ClientSecret"`
	RedirectURL  string `json:"redirect_url" flag:"redirect-url" reg:"RedirectURL"`
	TokenJson    string `json:"token_json" flag:"token-json" reg:"TokenJson"`
	AlbumFilter  string `json:"album_filter" flag:"album-filter" reg:"AlbumFilter"`
}

type TokenPersistence interface {
	// SaveToken saves the token to a file.
	SaveToken(token *oauth2.Token) error
	// LoadToken loads the token from a file.
	LoadToken() (*oauth2.Token, error)
}

func (c *Config) Authenticate(ctx context.Context, starturl func(string)) error {
	token, err := gpfs.Authenticate(ctx, c.ClientID, c.ClientSecret, c.RedirectURL, starturl)
	if err != nil {
		return err
	}
	if err := c.SaveToken(token); err != nil {
		return err
	}
	return nil
}

func (c *Config) SaveToken(token *oauth2.Token) error {
	tokenJson, err := json.Marshal(token)
	if err != nil {
		return err
	}
	c.TokenJson = string(tokenJson)
	return nil
}

func (c *Config) LoadToken() (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(c.TokenJson), &token); err != nil {
		return nil, err
	}
	return &token, nil
}

// ToFileSystem implements bindings.BindingConfig.
func (c *Config) ToFileSystem(zerolog.Logger) (afero.Fs, error) {
	oauthconfig := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
		Endpoint:     google.Endpoint,
	}
	token, err := c.LoadToken()
	if err != nil {
		return nil, err
	}
	httpclient := oauthconfig.Client(context.Background(), token)
	// TODO persist token when updated by the client
	afilter := strings.Split(c.AlbumFilter, ",")
	return gpfs.NewFs(httpclient, afilter)
}

// Validate implements bindings.BindingConfig.
func (c *Config) Validate() error {
	return nil
}
