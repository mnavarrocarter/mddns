package mddns

import (
	"context"
	"github.com/mnavarrocarter/mddns/debug"
	"github.com/mnavarrocarter/mddns/ip"
	"github.com/mnavarrocarter/mddns/ip/google"
	"github.com/mnavarrocarter/mddns/provider"
	_ "github.com/mnavarrocarter/mddns/provider/google"
	"github.com/mnavarrocarter/mddns/provider/system"
	"github.com/sirupsen/logrus"
	"net"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	ConfigFile string
	Watch      bool
	CacheFile  string
	Interval   time.Duration
	Debug      int
}

func (a *App) Run() {
	logger := logrus.New()
	logger.Level = logrus.Level(a.Debug)
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	debug.MutateDefaultClient(logger.Debugf)

	schedule := ip.Scheduler(a.Interval, a.CacheFile)
	gb := google.NewBackend(nil)

	updater, err := system.NewUpdater(a.ConfigFile, a.Watch, logger.Errorf)
	if err != nil {
		logger.Fatalf("provider error: %s", err.Error())
	}

	ips := generateIps(logger, schedule(gb))
	runLoop(ctx, logger, ips, updater)
}

func generateIps(logger *logrus.Logger, provider ip.Provider) <-chan net.IP {
	ipChan := make(chan net.IP)

	// This goroutine listens for new ip addresses and passes them to the ip channel
	go func() {
		logger.Infof("watching for ip changes")
		for {
			ctx := context.Background()

			i, err := provider.GetIP(ctx)
			if err != nil {
				logger.Errorf("ip provider error: %s", err.Error())
				continue
			}

			ipChan <- i
		}
	}()

	return ipChan
}

func runLoop(ctx context.Context, logger *logrus.Logger, ips <-chan net.IP, updater provider.Updater) {

	logger.Infof("starting the main loop")
	defer func() {
		logger.Infof("main loop stopped")
	}()

	// This loop checks for new ips or signals
	for {
		select {
		case _ = <-ctx.Done():
			logger.Warning("closing signal received")
			return
		case i := <-ips:
			logger.Infof("new ip found: %s", i.String())

			err := updater.Update(ctx, i)
			if err != nil {
				logger.Errorf("failed to update ip: %s", err.Error())
				continue
			}

			logger.Infof("new ip updated: %s", i.String())
		}
	}
}
