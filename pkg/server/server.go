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

	msgCh := make(chan string, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				msgCh <- err.Error()
			}
		}
		close(msgCh)
	}()

	log.I("server started")

	var err error = nil

	select {
	case m, ok := <-msgCh:
		if ok {
			err = fmt.Errorf(m)
		}
	case <-ctx.Done():
		_ = a.shutdown(3 * time.Second)
	}

	log.I("server stopping...")

	if m, ok := <-msgCh; ok {
		err = fmt.Errorf(m)
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
