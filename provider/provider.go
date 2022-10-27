package provider

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"sync"
)

var ErrUnsupported = errors.New("unsupported driver")

var mtx sync.Mutex
var drivers = make([]Driver, 0)

type Driver interface {
	Configure(uri *url.URL) (Updater, error)
}

type DriverFn func(uri *url.URL) (Updater, error)

func (u DriverFn) Configure(uri *url.URL) (Updater, error) {
	return u(uri)
}

type Updater interface {
	Update(ctx context.Context, ip net.IP) error
}

func Register(driver Driver) {
	mtx.Lock()
	defer mtx.Unlock()
	drivers = append(drivers, driver)
}

func Make(uri string) (Updater, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	for _, d := range drivers {
		u, err := d.Configure(u)
		if err == ErrUnsupported {
			continue
		}

		if err != nil {
			return nil, err
		}

		return u, nil
	}

	return nil, fmt.Errorf("%w: driver %q not registered", ErrUnsupported, u.Scheme)
}
