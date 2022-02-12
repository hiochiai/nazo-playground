package main

import (
	"context"
	"fmt"
	golog "log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hiochiai/nazo-playground/pkg/config"
	"github.com/hiochiai/nazo-playground/pkg/log"
	"github.com/hiochiai/nazo-playground/pkg/server"
)

func main() {

	ctx, cancel := contextWithSignal()
	defer cancel()

	err, exitCode := run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCode)
	}
}

func contextWithSignal() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	return ctx, cancel
}

func run(ctx context.Context) (error, int) {

	c, err := configure()
	if err != nil {
		return err, 2
	}

	initLog(c)

	nazoServer := server.NewNazoServer(c)
	if err := nazoServer.Serve(ctx); err != nil {
		return err, 1
	}

	return nil, 0
}

func initLog(c *config.Config) {
	log.Logger = golog.New(os.Stdout, "", golog.Ldate|golog.Ltime)
	switch {
	case strings.EqualFold(c.LogLevel, "debug"):
		log.Level = log.LevelD
	case strings.EqualFold(c.LogLevel, "info"):
		log.Level = log.LevelI
	case strings.EqualFold(c.LogLevel, "warn"):
		log.Level = log.LevelW
	case strings.EqualFold(c.LogLevel, "error"):
		log.Level = log.LevelE
	default:
		log.Level = log.LevelI
	}
}

func configure() (*config.Config, error) {

	c := config.NewDefaultConfig()
	if err := c.ConfigureWithEnv(); err != nil {
		return nil, err
	}
	if err := c.ConfigureWithArgs(); err != nil {
		return nil, err
	}

	return c, nil
}
