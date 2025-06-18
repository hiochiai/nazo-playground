package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hiochiai/nazo-playground/pkg/config"
	"github.com/hiochiai/nazo-playground/pkg/log"
)

type NazoServer struct {
	server http.Server
}

func NewNazoServer(c *config.Config) *NazoServer {

	return &NazoServer{
		server: http.Server{
			Addr:    fmt.Sprintf(":%d", c.Server.Port),
			Handler: NewHandler(c),
		},
	}
}

func (a *NazoServer) Serve(ctx context.Context) error {

	log.If("server starting with port %v", a.server.Addr)

	errs := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				errs <- err
			}
		}
		close(errs)
	}()

	log.I("server started")

	var err error

	select {
	case err = <-errs:
		log.E("server failed: " + err.Error())
	case <-ctx.Done():
		log.I("server stopping...")
		_ = a.shutdown(3 * time.Second)
		err = <-errs
	}

	log.I("server stopped")

	return err
}

func (a *NazoServer) shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
