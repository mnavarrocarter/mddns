package main

import (
	"flag"
	"github.com/mnavarrocarter/mddns"
	"os"
	"path/filepath"
	"time"
)

func main() {
	app := parseFlags()
	app.Run()
}

func parseFlags() *mddns.App {
	app := &mddns.App{}

	set := flag.NewFlagSet("mddns", flag.ExitOnError)
	set.StringVar(&app.ConfigFile, "config", "mddns.txt", "The config file where dnns entries are defined")
	set.BoolVar(&app.Watch, "watch", false, "Watches for changes to the config file")
	set.StringVar(&app.CacheFile, "cache", filepath.Join(os.TempDir(), "mddns_ip.txt"), "File where the current ip is stored")
	set.DurationVar(&app.Interval, "interval", time.Second*30, "The interval to check for new ip changes")
	set.IntVar(&app.Debug, "debug", 4, "Log output level")

	_ = set.Parse(os.Args[1:])

	return app
}
