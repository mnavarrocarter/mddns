package google

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mnavarrocarter/mddns/provider"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

const scheme = "google"
const updateUrl = "https://domains.google.com/nic/update"

var ErrNoHost = errors.New("missing or inexistent host")
var ErrBadAuth = errors.New("wrong authentication credentials")
var ErrNotFQDN = errors.New("host is not fully qualified domain name")
var ErrBadAgent = errors.New("bad user agent provided")
var ErrAbuse = errors.New("too many requests")
var Err911 = errors.New("unknown server error")
var ErrConflict = errors.New("conflicting A or AAAA record")

func init() {
	provider.Register(provider.DriverFn(NewUpdater))
}

func NewUpdater(uri *url.URL) (provider.Updater, error) {
	if uri.Scheme != scheme {
		return nil, provider.ErrUnsupported
	}

	pass, _ := uri.User.Password()

	return &updater{
		client:   http.DefaultClient,
		username: uri.User.Username(),
		password: pass,
		hostname: uri.Hostname(),
		email:    uri.Query().Get("email"),
	}, nil
}

type updater struct {
	client   *http.Client
	username string
	password string
	hostname string
	email    string
}

func (u *updater) Update(ctx context.Context, ip net.IP) error {
	query := url.Values{
		"hostname": []string{u.hostname},
		"myip":     []string{ip.String()},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, updateUrl+"?"+query.Encode(), http.NoBody)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+u.auth())
	req.Header.Add("User-Agent", u.userAgent())

	res, err := u.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return ErrBadAuth
	}

	defer func(c io.Closer) { _ = c.Close() }(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	status := string(b)

	if strings.Contains(status, "good") || strings.Contains(status, "nochg") {
		return nil
	}

	switch status {
	case "badauth":
		return ErrBadAuth
	case "nohost":
		return ErrNoHost
	case "notfqdn":
		return ErrNotFQDN
	case "badagent":
		return ErrBadAgent
	case "abuse":
		return ErrAbuse
	case "911":
		return Err911
	case "conflict A":
	case "conflict AAAA":
		return ErrConflict
	}

	return fmt.Errorf("unkndown error: %s", string(b))
}

func (u *updater) auth() string {
	auth := fmt.Sprintf("%s:%s", u.username, u.password)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (u *updater) userAgent() string {
	userAgent := "MDDNS v1"
	if u.email != "" {
		userAgent = " " + u.email
	}

	return userAgent
}
