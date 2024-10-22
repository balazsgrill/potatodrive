package server

import (
	"errors"
	"net/http"

	"github.com/spf13/afero"
)

type Config struct {
	KeyId     string `flag:"keyid,Access Key ID" reg:"KeyID" json:"keyid"`
	KeySecret string `flag:"secret,Access Key Secret" reg:"KeySecret" json:"secret"`
	Directory string `flag:"directory,Directory to serve" reg:"Directory" json:"directory"`
	Pattern   string `flag:"pattern,HTTP path pattern" reg:"Pattern" json:"pattern"`
}

func (c *Config) Validate() error {
	if c.KeyId == "" {
		return errors.New("keyid is mandatory")
	}
	if c.KeySecret == "" {
		return errors.New("secret is mandatory")
	}
	if c.Directory == "" {
		return errors.New("directory is mandatory")
	}
	if c.Pattern == "" {
		return errors.New("pattern is mandatory")
	}
	return nil
}

func (c *Config) ToHandler() (string, http.HandlerFunc, error) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), c.Directory)
	return c.Pattern, c.authMiddleware(Handler(fs)), nil
}

func (c *Config) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providedKeyID := r.Header.Get("X-Key-ID")
		providedSecret := r.Header.Get("X-Secret")

		if providedKeyID != c.KeyId || providedSecret != c.KeySecret {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
