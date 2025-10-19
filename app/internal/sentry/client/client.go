package client

import (
	"context"
	"io"
	"net/http"
	"path"
	"time"

	errorscath "catcher/pkg/errors"

	"github.com/cockroachdb/errors"
)

const method = "GET"

type Config struct {
	Url     string
	Org     string
	Token   string
	Timeout int
}

type Client struct {
	Config
}

func New(config Config) *Client {
	return &Client{
		Config: config}
}

func (c Client) Request(endpoint string) ([]byte, error) {

	const op = "sentry.client.request"

	url := c.Url + path.Join("organizations", c.Org, endpoint)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, op)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(c.Timeout)*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithMessage(errorscath.RequestRecovered(err), op)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("sentry respond with status: %d, %s", resp.StatusCode, op)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, op)
	}

	return body, nil
}
