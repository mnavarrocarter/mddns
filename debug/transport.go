package debug

import (
	"net/http"
)

type LogFn func(string, ...any)

func NewClient(log LogFn) *http.Client {
	client := *http.DefaultClient
	client.Transport = NewTransport(http.DefaultTransport, log)
	return &client
}

func MutateDefaultClient(log LogFn) func() {
	client := http.DefaultClient
	http.DefaultClient = NewClient(log)

	return func() {
		http.DefaultClient = client
	}
}

func NewTransport(next http.RoundTripper, log LogFn) *Transport {
	return &Transport{
		next: next,
		log:  log,
	}
}

type Transport struct {
	next http.RoundTripper
	log  LogFn
}

func (l *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := l.next.RoundTrip(req)
	if err != nil {
		return res, err
	}

	l.log("http %d response for %s %s", res.StatusCode, req.Method, req.URL.String())

	return res, nil
}
