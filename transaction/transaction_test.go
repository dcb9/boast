package transaction

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestTransport_RoundTrip(t *testing.T) {

}

func TestNewReq(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/welcome?user=bob", nil)
	rq := NewReq(req)
	t.Logf("%#v\n", rq)

	reqBody := bytes.NewReader([]byte(`{"name": "bob"}`))
	req, _ = http.NewRequest("POST", "http://localhost/welcome?user=bob", reqBody)
	rq = NewReq(req)
	t.Logf("%#v\n", rq)

	req, _ = http.NewRequest("PUT", "http://localhost/welcome?user=bob", reqBody)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8,zh-CN;q=0.6,zh;q=0.4")
	rq = NewReq(req)
	t.Logf("%#v\n", rq)
}

func TestServe(t *testing.T) {
	rp := NewReverseProxy(t)
	via = "1.1 " + string(rp.URL[7:])

	reqBody := bytes.NewReader([]byte(`{"name": "bob"}`))
	req, err := http.NewRequest("POST", rp.URL+"/welcome", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	{
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.8,zh-CN;q=0.6,zh;q=0.4")
		t.Log("Request BEGIAN", "====================")
		t.Log("MethodPath: ", req.Method+" "+req.URL.Path)
		t.Log("Headers:")
		for key, vals := range req.Header {
			t.Log("\t", key, " : ", vals)
		}

		t.Log("Request END", "=======================")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	{
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		t.Log("Response BEGIAN", "====================")
		t.Log("Response Status: ", resp.Status)
		t.Log("Response Body: ", string(body))
		t.Log("Response Headers: ")

		rawHeaders := make([]string, 0, len(resp.Header))
		for key, vals := range resp.Header {
			t.Log("\t", key, " : ", vals)
			rawHeaders = append(rawHeaders, key+": "+strings.Join(vals, ",")+"\r\n")
		}

		t.Log("RAW Response:")
		rawResp := fmt.Sprintf(
			"%s %s\r\n%s\r\n%s",
			resp.Proto, resp.Status,
			strings.Join(rawHeaders, ""), body,
		)
		t.Log(rawResp)
		t.Log("Response END", "=======================")
	}
}

func NewReverseProxy(t *testing.T) *httptest.Server {
	bs := NewBackendServer(t)
	target, err := url.Parse(bs.URL)
	if err != nil {
		t.Error(err)
	}

	proxy := NewSingleHostReverseProxy(target)
	proxy.Transport = transport

	return httptest.NewServer(proxy)
}

func NewBackendServer(t *testing.T) *httptest.Server {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bs, _ := ioutil.ReadAll(r.Body)
		t.Log("GOT Request URL: ", string(r.URL.Host))
		t.Log("GOT Request Body: ", string(bs))
		body := []byte("Welcome to Boast!")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.Write(body)
	}))
	return backendServer
}
