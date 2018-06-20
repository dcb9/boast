package transaction

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Req struct {
	URL    *url.URL    `json:"URL"`
	Method string      `json:"Method"`
	Proto  string      `json:"Proto"`
	Header http.Header `json:"Header"`
	Body   []byte      `json:"Body"`
}

type Resp struct {
	Proto  string
	Header http.Header `json:"Header"`
	Body   []byte      `json:"Body"`
	Status string      `json:"Status"`
}

type Ts struct {
	ID         uuid.UUID     `json:"ID"`
	Req        *Req          `json:"Req"`
	Resp       *Resp         `json:"Resp"`
	ClientAddr string        `json:"ClientAddr"`
	BeginAt    time.Time     `json:"BeginAt"`
	EndAt      time.Time     `json:"EndAt"`
}

func NewReq(req *http.Request) *Req {
	rURL := new(url.URL)
	rawReqURL := req.Context().Value("rawRequestURL")
	if rawReqURL != nil {
		*rURL = *rawReqURL.(*url.URL)
	} else {
		*rURL = *req.URL
	}

	var body []byte
	req.Body, body = copyBody(req.Body)

	return &Req{
		URL:    rURL,
		Method: req.Method,
		Proto:  req.Proto,
		Header: CopyHeader(req.Header),
		Body:   body,
	}
}

func NewResp(resp *http.Response) *Resp {
	var body []byte
	resp.Body, body = copyBody(resp.Body)

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		g, _ := gzip.NewReader(bytes.NewReader(body))
		body, _ = ioutil.ReadAll(g)
	}

	return &Resp{
		Proto:  resp.Proto,
		Header: CopyHeader(resp.Header),
		Body:   body,
		Status: resp.Status,
	}
}

func CopyHeader(h http.Header) http.Header {
	header := make(http.Header)
	for k, v := range h {
		header[k] = v
	}

	return header
}

func copyBody(r io.ReadCloser) (io.ReadCloser, []byte) {
	body := make([]byte, 0)
	if r != nil {
		body, _ = ioutil.ReadAll(r)
		r = ioutil.NopCloser(bytes.NewReader(body))
	}
	return r, body
}
