package main

import (
	"flag"
	"github.com/mnavarrocarter/mddns"
	"os"
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
	set.DurationVar(&app.Interval, "interval", time.Second*30, "The interval to check for new ip changes")
	set.IntVar(&app.Debug, "debug", 4, "Log output level")

	_ = set.Parse(os.Args[1:])

	return app
}
