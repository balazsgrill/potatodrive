package client

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type Config struct {
	URL       string `flag:"url,Proxy URL" reg:"URL"`
	KeyId     string `flag:"keyid,Access Key ID" reg:"KeyID"`
	KeySecret string `flag:"secret,Access Key Secret" reg:"KeySecret"`
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return errors.New("url is mandatory")
	}
	if c.KeyId == "" {
		return errors.New("keyid is mandatory")
	}
	if c.KeySecret == "" {
		return errors.New("secret is mandatory")
	}
	return nil
}

type authenticator struct {
	*Config
	delegate http.RoundTripper
}

func (a *authenticator) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("X-Key-ID", a.KeyId)
	r.Header.Add("X-Secret", a.KeySecret)
	return a.delegate.RoundTrip(r)
}

func (c *Config) ToFileSystem(logger zerolog.Logger) (afero.Fs, error) {
	httpclient := &http.Client{}
	httpclient.Transport = &authenticator{
		Config:   c,
		delegate: http.DefaultTransport,
	}
	return Connect(c.URL, httpclient)
}
