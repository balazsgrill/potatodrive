package gpfs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/balazsgrill/potatodrive/assets"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authprocess struct {
	oauthConfig  *oauth2.Config
	opts         []oauth2.AuthCodeOption
	state        string
	token        chan *oauth2.Token
	errorchannel chan error
	server       *http.Server
}

func Authenticate(ctx context.Context, clientID, clientSecret, redirectURL string, starturl func(string)) (*oauth2.Token, error) {
	authprocess := &authprocess{
		token:        make(chan *oauth2.Token),
		errorchannel: make(chan error),
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary"},
			Endpoint:     google.Endpoint,
		},
		opts: []oauth2.AuthCodeOption{
			// offline access is required to get a refresh token
			// and prompt=consent is required to get a refresh token on every login
			oauth2.SetAuthURLParam("access_type", "offline"),
			oauth2.SetAuthURLParam("prompt", "consent"),
		},
		state: uuid.New().String(),
	}

	u, err := url.Parse(authprocess.oauthConfig.RedirectURL)
	if err != nil {
		log.Fatalf("Invalid redirect URL: %v", err)
	}
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else if u.Scheme == "http" {
			port = "80"
		} else {
			log.Fatalf("Unknown scheme: %s", u.Scheme)
		}
		log.Fatal("No port specified in redirect URL")
	}
	callbackPath := "/callback"
	if u.Path != "" && u.Path != "/" {
		callbackPath = u.Path
	}
	baseurl := fmt.Sprintf("%s://%s:%s", u.Scheme, u.Hostname(), port)
	fmt.Printf("Starting server at %s\n", baseurl)
	servemux := http.NewServeMux()
	servemux.HandleFunc("/", authprocess.handleLogin)
	servemux.HandleFunc(callbackPath, authprocess.handleCallback)
	authprocess.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: servemux,
	}
	go func() {
		err := authprocess.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			authprocess.errorchannel <- err
		}
	}()
	defer authprocess.server.Shutdown(context.Background())
	defer close(authprocess.token)
	starturl(baseurl)
	select {
	case <-ctx.Done():
		authprocess.server.Shutdown(context.Background())
		return nil, fmt.Errorf("context cancelled")
	case token := <-authprocess.token:
		return token, nil
	case err := <-authprocess.errorchannel:
		if err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("no token received")
}

func (ap *authprocess) handleLogin(w http.ResponseWriter, r *http.Request) {
	url := ap.oauthConfig.AuthCodeURL(ap.state, ap.opts...)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (ap *authprocess) handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != ap.state {
		http.Error(w, "State is invalid", http.StatusBadRequest)
		ap.errorchannel <- fmt.Errorf("state is invalid")
		return
	}

	code := r.URL.Query().Get("code")
	token, err := ap.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		ap.errorchannel <- fmt.Errorf("state is invalid")
		return
	}
	ap.token <- token

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write(assets.AuthSuccessHtml)
}
