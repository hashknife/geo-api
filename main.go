package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/hashknife/common/services"
	"github.com/hashknife/geo-api/bindings"
	"github.com/hashknife/geo-api/config"
)

var gitSHA string
var signalsChan = make(chan os.Signal, 1)

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}

func main() {
	var configFlag string
	flag.StringVar(&configFlag, "c", "", "")
	flag.Parse()
	if configFlag == "" {
		fmt.Println("config file required")
		os.Exit(1)
	}

	logger := kitlog.NewJSONLogger(os.Stdout)
	rand.Seed(time.Now().UnixNano())
	root := context.Background()

	// errChan sends termination signals and blocks main goroutine
	errChan := make(chan error)

	conf, err := config.Load(configFlag, logger)
	if err != nil {
		logger.Log("fatal", err.Error())
		os.Exit(1)
	}

	hlp := &bindings.HTTPListenerParams{
		Logger:  logger,
		Root:    root,
		ErrChan: errChan,
		Config:  conf,
	}
	t38 := services.NewTile38(*conf.Tile38.Hostname)
	bindings.StartHealthCheckHTTPListener(hlp, gitSHA)
	bindings.StartApplicationHTTPListener(hlp, t38)

	logger.Log("fatal", <-errChan)
}
