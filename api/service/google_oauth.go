package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthService struct {
	config *oauth2.Config
}

type GoogleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewGoogleOAuthService(clientID, clientSecret, redirectURL string) *GoogleOAuthService {
	return &GoogleOAuthService{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// Enabled reports whether Google OAuth is fully configured. RedirectURL is
// included so a partially-configured service (one of the three env vars
// missing) is treated as disabled instead of producing a flow that fails
// late at Google's redirect_uri_mismatch check.
func (s *GoogleOAuthService) Enabled() bool {
	return s.config.ClientID != "" && s.config.ClientSecret != "" && s.config.RedirectURL != ""
}

// AuthURL returns the URL the browser should be redirected to in order to start
// the OAuth flow. The state value should be persisted (e.g. as a short-lived
// cookie) and verified on the callback to prevent CSRF.
func (s *GoogleOAuthService) AuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// Exchange swaps the authorization code for an access token and fetches the
// user profile from Google.
func (s *GoogleOAuthService) Exchange(ctx context.Context, code string) (*GoogleUserInfo, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth exchange failed: %w", err)
	}

	client := s.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch userinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info GoogleUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo: %w", err)
	}
	if info.Sub == "" || info.Email == "" {
		return nil, errors.New("incomplete user info from google")
	}
	return &info, nil
}

// GenerateState produces a cryptographically random opaque string used as the
// OAuth state parameter.
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
