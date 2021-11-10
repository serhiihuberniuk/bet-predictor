package fetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/serhiihuberniuk/bet-predictor/fetcher/es-fetcher/es-models"
)

type ESFetcher struct {
	client http.Client
	token  string
}

func NewESFetcher(ctx context.Context, apiKey string) (*ESFetcher, error) {
	client := http.Client{}

	token, err := getToken(ctx, client, apiKey)
	if err != nil {
		return nil, fmt.Errorf("error while getting access token: %w", err)
	}

	return &ESFetcher{
		client: client,
		token:  token,
	}, nil
}

func getToken(ctx context.Context, c http.Client, apiKey string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, "https://oauth2.elenasport.io/oauth2/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("error while creating request: %w", err)
	}

	req.Header.Add("Authorization", "Basic "+apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return "", fmt.Errorf("error while sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("response with status: " + resp.Status)
	}

	var authResponse esmodels.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return "", fmt.Errorf("error while decoding response: %w", err)
	}

	return authResponse.AccessToken, nil
}
