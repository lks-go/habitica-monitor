// Package habitica — клиент внешнего Habitica API (адаптер выходного порта).
package habitica

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://habitica.com/api/v3"

// Stats — подмножество ответа GET /user.data.stats, нужное монитору.
type Stats struct {
	HP  float64 `json:"hp"`
	MP  float64 `json:"mp"`
	Exp float64 `json:"exp"`
	GP  float64 `json:"gp"`
	Lvl int     `json:"lvl"`
}

// Client вызывает Habitica API.
type Client struct {
	httpClient *http.Client
	xClient    string // формат "<your-user-id>-<appname>", требование Habitica
}

// NewClient создаёт клиент. xClient обязателен в заголовке каждого запроса.
func NewClient(xClient string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		xClient:    xClient,
	}
}

type userResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Stats Stats `json:"stats"`
	} `json:"data"`
}

// GetUserStats делает GET /user от имени пользователя (по его id и apiKey)
// и возвращает текущие статы.
func (c *Client) GetUserStats(ctx context.Context, apiUser, apiKey string) (Stats, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		baseURL+"/user?userFields=stats", nil)
	if err != nil {
		return Stats{}, err
	}
	req.Header.Set("x-api-user", apiUser)
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("x-client", c.xClient)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Stats{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return Stats{}, fmt.Errorf("habitica: status %d: %s", resp.StatusCode, string(b))
	}

	var ur userResponse
	if err := json.NewDecoder(resp.Body).Decode(&ur); err != nil {
		return Stats{}, fmt.Errorf("habitica: decode: %w", err)
	}
	return ur.Data.Stats, nil
}
