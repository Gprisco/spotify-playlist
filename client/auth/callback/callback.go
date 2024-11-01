package callback

import (
	"net/http"
	"time"
)

type CallbackResult struct {
	Code string
	Err  string
}

type CallbackContext struct {
	channel chan *CallbackResult
}

type CallbackHandler func(timeout time.Duration) *CallbackResult

func (p *CallbackContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Asynchronously send the code query param back (and error if present)
		code := r.URL.Query().Get("code")
		err := r.URL.Query().Get("error")
		p.channel <- &CallbackResult{Code: code, Err: err}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	}
}
