package oauth

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// Config example (fill these with your values)
var OAuthConfig = &oauth2.Config{
	ClientID:     "YOUR_FACEBOOK_APP_ID",
	ClientSecret: "YOUR_FACEBOOK_APP_SECRET",
	RedirectURL:  "http://localhost:8080/oauth/callback",
	Scopes:       []string{"email", "public_profile"},
	Endpoint:     facebook.Endpoint,
}

// ConnectToFacebook exchanges a code for an access token and returns a message or error
func ConnectToFacebook(code string) (string, error) {
	ctx := context.Background()

	// Exchange code for access token
	token, err := OAuthConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Use token to fetch user info from Facebook
	client := OAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email")
	if err != nil {
		return "", fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Facebook API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Return user info as string message
	return string(body), nil
}
