package fetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/serhiihuberniuk/bet-predictor/fetcher/es_fetcher/es_models"
)

type ESFetcher struct {
	client http.Client
	apiKey string
	token  string
}

func NewESFetcher(ctx context.Context, apiKey string) (*ESFetcher, error) {
	fetcher := &ESFetcher{
		client: http.Client{},
		apiKey: apiKey,
	}

	if err := fetcher.getToken(ctx); err != nil {
		return nil, fmt.Errorf("error while getting access token: %w", err)
	}

	return fetcher, nil
}

func (f *ESFetcher) getToken(ctx context.Context) error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, "https://oauth2.elenasport.io/oauth2/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error while creating request: %w", err)
	}

	req.Header.Add("Authorization", "Basic "+f.apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var authResponse esmodels.AuthResponse
	decodeFn := func(reader io.Reader) error {
		if err := json.NewDecoder(reader).Decode(&authResponse); err != nil {
			return fmt.Errorf("error while decoding response: %w", err)
		}

		return nil
	}

	if err = f.doRequest(ctx, req, decodeFn); err != nil {
		return fmt.Errorf("error while doing request to get access token :%w", err)
	}

	f.token = authResponse.AccessToken

	return nil
}

func (f *ESFetcher) doRequest(ctx context.Context, req *http.Request, decodeFn func(reader io.Reader) error) error {
	if req.Header.Get("Authorization") == "" {
		req.Header.Add("Authorization", "Bearer "+f.token)
	}

	resp, err := f.client.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error while sending request: %w", err)
	}

	if resp.StatusCode == http.StatusForbidden {
		if err = f.getToken(ctx); err != nil {
			return fmt.Errorf("error while getting token: %w", err)
		}

		return f.doRequest(ctx, req, decodeFn)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("response with status: " + resp.Status)
	}

	if err := decodeFn(resp.Body); err != nil {
		if err != nil {
			return fmt.Errorf("error while decoding data: %w", err)
		}
	}

	return nil
}
