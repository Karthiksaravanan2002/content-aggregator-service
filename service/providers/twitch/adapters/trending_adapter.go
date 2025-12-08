package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
)

var (
	cachedToken     string
	tokenExpiryTime time.Time
	tokenMu         sync.Mutex
)

type twitchTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func getAppToken(cfg props.TwitchConfig, logger *zap.Logger) (string, error) {
	tokenMu.Lock()
	defer tokenMu.Unlock()

	if cachedToken != "" && time.Now().Before(tokenExpiryTime) {
		return cachedToken, nil
	}

	url := "https://id.twitch.tv/oauth2/token?client_id=" + cfg.ClientID +
		"&client_secret=" + cfg.ClientSecret +
		"&grant_type=client_credentials"

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tok twitchTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return "", err
	}

	cachedToken = tok.AccessToken
	tokenExpiryTime = time.Now().Add(time.Duration(tok.ExpiresIn-60) * time.Second)
	return cachedToken, nil
}

func CallTrendingStreams(cfg props.TwitchConfig, logger *zap.Logger) func(context.Context, *models.ProviderRequest) ([]*models.ContentItem, errors.AppError) {
	return func(ctx context.Context, req *models.ProviderRequest) ([]*models.ContentItem, errors.AppError) {

		token, err := getAppToken(cfg, logger)
		if err != nil {
			return nil, errors.ProviderError(http.StatusBadGateway, err)
		}

		url := "https://api.twitch.tv/helix/streams?first=20"

		httpReq, _ := http.NewRequest("GET", url, nil)
		httpReq.Header.Set("Client-ID", cfg.ClientID)
		httpReq.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			return nil, errors.ProviderError(http.StatusBadGateway, err)
		}
		defer resp.Body.Close()

		var raw struct {
			Data []map[string]any `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&raw)

		items := make([]*models.ContentItem, 0, len(raw.Data))
		for _, stream := range raw.Data {
			items = append(items, MapTwitchStream(stream))
		}

		return items, nil
	}
}
