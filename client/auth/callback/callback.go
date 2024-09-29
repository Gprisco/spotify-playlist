package callback

import "net/http"

type CallbackResult struct {
	Code string
	Err  string
}

type CallbackHandler struct {
	channel chan *CallbackResult
}

func (p *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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