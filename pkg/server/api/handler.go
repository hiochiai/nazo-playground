package api

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/hiochiai/nazo-playground/pkg/config"
	"github.com/hiochiai/nazo-playground/pkg/log"
)

const PathPrefix = "/api/v1/nazo"

type Handler struct {
	handleFuncs map[string]func(w http.ResponseWriter, r *http.Request)
}

type NazoPage struct {
	Id       string
	NextId   string
	Answer   string
	Contents string
}

func NewHandler(pages []NazoPage) *Handler {

	funcs := make(map[string]func(w http.ResponseWriter, r *http.Request))

	for i, page := range pages {
		path := filepath.Join(PathPrefix, page.Id)
		funcs[path] = makeGetPageHandler(page.Contents)
		if i == 0 {
			funcs[PathPrefix] = funcs[path]     // Redirect to first page
			funcs[PathPrefix+`/`] = funcs[path] // Redirect to first page
		}

		if len(page.Answer) == 0 {
			continue
		}

		path = filepath.Join(path, `answer`)
		funcs[path] = makePostAnswerHandler(page.Answer, page.NextId)
		if i == 0 {
			funcs[filepath.Join(PathPrefix, `answer`)] = funcs[path] // Redirect to first page
		}
	}

	return &Handler{
		handleFuncs: funcs,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleFunc, ok := h.handleFuncs[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	handleFunc(w, r)
}

func makeGetPageHandler(contents string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		contents, _ = config.EvalPagesContents(contents)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(contents))
	}
}

func makePostAnswerHandler(correctAnswer string, next string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		req, err := GetJsonRequestBody(r)
		if err != nil {
			log.If(`failed to parse request body: %v`, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		answer, ok := req["answer"].(string)
		if !ok {
			log.If(`failed to parse request body for answer`)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var response string
		if answer == correctAnswer {
			response = fmt.Sprintf(`{"result":%v,"next":"%s"}`, true, next)
		} else {
			response = fmt.Sprintf(`{"result":%v}`, false)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(response))
	}
}
