package system

import (
	"bufio"
	"context"
	"errors"
	"github.com/mnavarrocarter/mddns/provider"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

type LogFn func(string, ...any)

var ErrMultiple = errors.New("multiple errors while updating")

func NewUpdater(file string, watch bool, log LogFn) (*Updater, error) {
	var err error

	updater := &Updater{
		mtx:  sync.Mutex{},
		file: file,
		log:  log,
	}

	updater.updaters, err = parseUpdaters(file)
	if err != nil {
		return nil, err
	}

	return updater, nil
}

type Updater struct {
	mtx      sync.Mutex
	file     string
	updaters []provider.Updater
	log      LogFn
}

func (u *Updater) Update(ctx context.Context, ip net.IP) error {
	wg := sync.WaitGroup{}

	errs := make([]error, 0)

	for i := 0; i <= len(u.updaters)-1; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			err := u.updaters[i].Update(ctx, ip)
			if err != nil {
				u.log("update error: %s", err.Error())
				errs = append(errs, err)
			}
		}(i)
	}

	wg.Wait()
	if len(errs) > 0 {
		return ErrMultiple
	}

	return nil
}

func parseUpdaters(file string) ([]provider.Updater, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer func(c io.Closer) { _ = f.Close() }(f)

	updaters := make([]provider.Updater, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()

		// Empty line, we ignore
		if text == "" {
			continue
		}

		// Comment line, we ignore
		if strings.Index(text, "#") == 0 {
			continue
		}

		p, err := provider.Make(text)
		if err != nil {
			return nil, err
		}

		updaters = append(updaters, p)
	}

	return updaters, nil
}
