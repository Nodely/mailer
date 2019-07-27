package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/nodely/notify/pkg/receiver"
	"github.com/nodely/notify/pkg/sender"
	"github.com/op/go-logging"
)

func init() {
	var format = logging.MustStringFormatter(
		`%{color} ▶ %{level:-8s}%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
}

type options struct {
	Port int    `short:"p" long:"port" default:"8080" env:"NOTIFY_APP_PORT" description:"Port for Receiver Application"`
	Mode string `long:"mode" default:"mixed" env:"NOTIFY_APP_MODE" description:"Application mode: mixed, receiver, sender"`
}

func main() {

	log := logging.MustGetLogger("notify")

	var opts options

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warning("Interrupt signal")
		cancel()
	}()

	log.Infof("Application Mode: %s", opts.Mode)

	// we need to launch http server only for mixed or receiver modes
	if opts.Mode != "sender" {
		app := receiver.New(ctx, opts.Port)
		app.Launch()
	}

	// we need to listen queue only for mixed or sender modes
	if opts.Mode != "receiver" {
		app := sender.New(ctx)
		app.Launch()
	}

}
