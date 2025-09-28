package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"line-oa-backend/config"
)

type LINEOAuthService struct {
	config *config.Config
}

type LINETokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type LINEProfile struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

type LINEErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewLINEOAuthService(cfg *config.Config) *LINEOAuthService {
	return &LINEOAuthService{
		config: cfg,
	}
}

// GetAuthURL generates the LINE OAuth authorization URL
func (s *LINEOAuthService) GetAuthURL(state string) string {
	baseURL := "https://access.line.me/oauth2/v2.1/authorize"
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", s.config.LINEChannelID)
	params.Add("redirect_uri", s.config.LINERedirectURI)
	params.Add("state", state)
	params.Add("scope", "profile openid")

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// ExchangeCodeForToken exchanges authorization code for access token
func (s *LINEOAuthService) ExchangeCodeForToken(code string) (*LINETokenResponse, error) {
	tokenURL := "https://api.line.me/oauth2/v2.1/token"
	
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.config.LINERedirectURI)
	data.Set("client_id", s.config.LINEChannelID)
	data.Set("client_secret", s.config.LINEChannelSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp LINEErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("failed to parse error response: %w", err)
		}
		return nil, fmt.Errorf("LINE API error: %s - %s", errorResp.Error, errorResp.ErrorDescription)
	}

	var tokenResp LINETokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// GetProfile gets user profile using access token
func (s *LINEOAuthService) GetProfile(accessToken string) (*LINEProfile, error) {
	profileURL := "https://api.line.me/v2/profile"

	req, err := http.NewRequest("GET", profileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get profile: status %d, body: %s", resp.StatusCode, string(body))
	}

	var profile LINEProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse profile response: %w", err)
	}

	return &profile, nil
}
