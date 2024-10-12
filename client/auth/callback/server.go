package callback

import (
	"context"
	"net/http"
	"time"
)

func HandleCallback(timeout time.Duration) *CallbackResult {
	channel := make(chan *CallbackResult)
	handler := &CallbackContext{channel: channel}

	// Spin up a server
	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	go server.ListenAndServe()

	// Catch the result
	result := <-channel

	// Shutdown the server right after receiving the callback result
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	server.Shutdown(ctx)

	return result
}
