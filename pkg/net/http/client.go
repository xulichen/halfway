package http

import (
	"go.elastic.co/apm/module/apmhttp"
	"net/http"
	"time"
)

var Client = NewClient(0)

// 接入apm的client
func NewClient(timeout time.Duration) *http.Client {
	// apmhttp.WithClientTrace() will show the detail of client trace including dns connect request tls response
	client := apmhttp.WrapClient(http.DefaultClient)
	client.Timeout = timeout
	return client
}
