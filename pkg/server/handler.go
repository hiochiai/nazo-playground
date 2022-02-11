package server

import (
	"net/http"
	"strings"

	"github.com/hiochiai/nazo-playground/pkg/config"
	"github.com/hiochiai/nazo-playground/pkg/server/api"
	"github.com/hiochiai/nazo-playground/pkg/server/static"
)

type Handler struct {
	mux map[string]http.Handler
}

func NewHandler(c *config.Config) http.Handler {

	l := len(c.Pages)
	pages := make([]api.NazoPage, 0)
	for i := 0; i < l; i++ {
		var next string
		if i+1 < l {
			next = c.Pages[i+1].Id
		}

		pages = append(pages,
			api.NazoPage{
				Id:       c.Pages[i].Id,
				NextId:   next,
				Answer:   c.Pages[i].Answer,
				Contents: c.Pages[i].Contents,
			})
	}

	staticDirPath := config.MakeConfStaticDirPath(c.ConfDirPath)

	h := Handler{
		mux: make(map[string]http.Handler),
	}
	h.mux[api.PathPrefix] = api.NewHandler(pages)
	h.mux[static.PathPrefix] = static.NewHandler(staticDirPath)

	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range h.mux {
		if strings.HasPrefix(r.URL.Path, k) {
			v.ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}
