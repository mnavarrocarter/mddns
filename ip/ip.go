package ip

import (
	"context"
	"net"
)

type Provider interface {
	GetIP(ctx context.Context) (net.IP, error)
}

type ProviderFn func(ctx context.Context) (net.IP, error)

func (i ProviderFn) GetIP(ctx context.Context) (net.IP, error) {
	return i(ctx)
}
