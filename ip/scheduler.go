package ip

import (
	"context"
	"net"
	"os"
	"time"
)

func Scheduler(d time.Duration, file string) func(next Provider) Provider {
	current := readCachedIp(file)

	return func(next Provider) Provider {
		return ProviderFn(func(ctx context.Context) (net.IP, error) {
			for {
				ip, err := next.GetIP(ctx)
				if err != nil {
					return nil, err
				}

				if !ip.Equal(current) {
					current = ip
					writeCachedIp(file, current)
					return current, nil
				}

				time.Sleep(d)
			}
		})
	}
}

func readCachedIp(file string) net.IP {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	ip := net.ParseIP(string(b))
	if ip == nil {
		return net.IPv4zero
	}

	return ip
}

func writeCachedIp(file string, ip net.IP) {
	_ = os.WriteFile(file, []byte(ip.String()), 0644)
}
