package google_test

import (
	"context"
	"github.com/mnavarrocarter/mddns/provider"
	"github.com/mnavarrocarter/mddns/provider/google"
	"net"
	"net/url"
	"testing"
)

func mustParseUrl(uri string) *url.URL {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	return u
}

func Test_Updater_ErrUnsupported(t *testing.T) {
	uri := mustParseUrl("google://")

	_, err := google.NewUpdater(uri)
	if err != provider.ErrUnsupported {

	}
}

func Test_Updater_Update(t *testing.T) {
	uri := mustParseUrl("google://user:pass@sub.domain.com?email=hello")

	updater, err := google.NewUpdater(uri)
	if err != nil {
		t.Fatal(err)
	}

	err = updater.Update(context.Background(), net.ParseIP("192.168.1.1"))
	if err != nil {
		t.Fatal(err)
	}
}
