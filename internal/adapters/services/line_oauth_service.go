package services

import (
	"encoding/json"
	"fmt"
	"io"
	"line-oa-backend/internal/domain/valueobjects"
	"net/http"
	"net/url"
	"strings"
)

// LINEOAuthServiceAdapter implements LINEOAuthService
type LINEOAuthServiceAdapter struct {
	channelID     string
	channelSecret string
	redirectURI   string
}

// NewLINEOAuthServiceAdapter creates a new LINE OAuth service adapter
func NewLINEOAuthServiceAdapter(channelID, channelSecret, redirectURI string) *LINEOAuthServiceAdapter {
	return &LINEOAuthServiceAdapter{
		channelID:     channelID,
		channelSecret: channelSecret,
		redirectURI:   redirectURI,
	}
}

// GetAuthURL generates LINE OAuth authorization URL
func (s *LINEOAuthServiceAdapter) GetAuthURL(state string) string {
	baseURL := "https://access.line.me/oauth2/v2.1/authorize"
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", s.channelID)
	params.Add("redirect_uri", s.redirectURI)
	params.Add("state", state)
	params.Add("scope", "profile openid")

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// ExchangeCodeForToken exchanges authorization code for access token
func (s *LINEOAuthServiceAdapter) ExchangeCodeForToken(code string) (*valueobjects.AuthToken, error) {
	tokenURL := "https://api.line.me/oauth2/v2.1/token"
	
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.redirectURI)
	data.Set("client_id", s.channelID)
	data.Set("client_secret", s.channelSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LINE token exchange failed: %s", string(body))
	}

	var tokenResp valueobjects.AuthToken
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// GetProfile retrieves user profile from LINE
func (s *LINEOAuthServiceAdapter) GetProfile(accessToken string) (*valueobjects.LINEProfile, error) {
	profileURL := "https://api.line.me/v2/profile"

	req, err := http.NewRequest("GET", profileURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LINE profile fetch failed: %s", string(body))
	}

	var profile valueobjects.LINEProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
