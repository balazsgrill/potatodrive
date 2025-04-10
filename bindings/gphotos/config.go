package gphotos

import (
	"context"

	"github.com/balazsgrill/potatodrive/gpfs"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	ClientID     string `json:"client_id" flag:"client-id" reg:"client_id"`
	ClientSecret string `json:"client_secret" flag:"client-secret" reg:"client_secret"`
	RedirectURL  string `json:"redirect_url" flag:"redirect-url" reg:"redirect_url"`
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
	httpclient := oauthconfig.Client(context.Background(), &oauth2.Token{})
	gpclient, err := gphotos.NewClient(httpclient)
	if err != nil {
		return nil, err
	}
	return gpfs.NewFs(gpclient), nil
}

// Validate implements bindings.BindingConfig.
func (c *Config) Validate() error {
	return nil
}
