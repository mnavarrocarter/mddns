package google

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
)

const ipCheckUrl = "https://domains.google.com/checkip"

func NewBackend(client *http.Client) *Backend {
	if client == nil {
		client = http.DefaultClient
	}

	return &Backend{client: client}
}

type Backend struct {
	client *http.Client
}

var ErrInvalidIP = errors.New("invalid ip")
var ErrInvalidResponse = errors.New("invalid response")

func (g *Backend) GetIP(ctx context.Context) (net.IP, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ipCheckUrl, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(c io.Closer) { _ = c.Close() }(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, ErrInvalidResponse
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(b))
	if ip == nil {
		return nil, ErrInvalidIP
	}

	return ip, nil
}
