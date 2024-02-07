package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/marcelluseasley/redditracker/config"
)

const (
	baseURL = "https://www.reddit.com"
)

type RedditToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenClient struct {
	httpClient *http.Client
}

type JWTTransport struct {
	AccessToken string
	Expiry      time.Time
	mu          sync.Mutex
	Transport   http.RoundTripper
	Config      *config.Config
}

func NewJWTTransport(accessToken string, expiry time.Time, config *config.Config) *JWTTransport {
	return &JWTTransport{
		AccessToken: accessToken,
		Expiry:      expiry,
		Transport:   http.DefaultTransport,
		Config:      config,
	}
}

func (t *JWTTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Expiry.Before(time.Now()) {
		tokenClient := NewTokenClient()

		newToken, err := tokenClient.GetToken(t.Config)
		if err != nil {
			return nil, err
		}

		t.AccessToken = newToken.AccessToken
		t.Expiry = ExpiresInToExpiry(newToken.ExpiresIn)
	}

	clonedRequest := cloneRequest(r)
	clonedRequest.Header.Add("Authorization", "Bearer "+t.AccessToken)
	return t.Transport.RoundTrip(clonedRequest)
}

func cloneRequest(r *http.Request) *http.Request {
	clonedRequest := r.Clone(r.Context())

	if r.Body != nil {
		clonedRequest.Body = r.Body
	}
	return clonedRequest
}

func NewTokenClient() *TokenClient {
	return &TokenClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (tc *TokenClient) GetToken(conf *config.Config) (*RedditToken, error) {
	requestData := url.Values{}
	requestData.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", baseURL+"/api/v1/access_token", strings.NewReader(requestData.Encode()))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(conf.RedditClientID, conf.RedditClientSecret)
	req.Header.Add("User-Agent", conf.UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := tc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tokenResp := &RedditToken{}
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, err
	}

	return tokenResp, nil
}
