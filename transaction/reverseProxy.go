package transaction

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"context"
)

var via string

type Transport struct {
	http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	beginAt := time.Now()
	tsReq := NewReq(req)
	if _, ok := req.Header["Accept-Encoding"]; ok {
		req.Header.Set("Accept-Encoding", "gzip")
	}

	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	endAt := time.Now()

	if via != "" {
		if v := resp.Header.Get("Via"); v == "" {
			resp.Header.Set("Via", via)
		} else {
			resp.Header.Set("Via", v+", "+via)
		}
	}

	{
		id, _ := uuid.NewUUID()

		t := Ts{
			ID:         id,
			Req:        tsReq,
			Resp:       NewResp(resp),
			ClientAddr: req.RemoteAddr,
			BeginAt:    beginAt,
			EndAt:      endAt,
		}

		go TsHub.Add(t)
	}

	return resp, nil
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		rURL := new(url.URL)
		*rURL = *req.URL
		rURL.Host = req.Host
		if rURL.Scheme == "" {
			rURL.Scheme = "http"
		}
		ctx := context.WithValue(req.Context(), "rawRequestURL", rURL)
		*req = *req.WithContext(ctx)

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
