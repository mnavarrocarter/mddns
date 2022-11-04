package ip

import (
	"context"
	"net"
	"time"
)

func Scheduler(d time.Duration) func(next Provider) Provider {
	current := net.IPv4zero

	return func(next Provider) Provider {
		return ProviderFn(func(ctx context.Context) (net.IP, error) {
			for {
				ip, err := next.GetIP(ctx)
				if err != nil {
					return nil, err
				}

				if !ip.Equal(current) {
					current = ip
					return current, nil
				}

				time.Sleep(d)
			}
		})
	}
}
