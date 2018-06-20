package transaction

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/dcb9/boast/config"
	"github.com/google/uuid"
)

var TsHub = &Hub{
	Transactions: make(map[uuid.UUID]*Ts),
	SortID:       make([]uuid.UUID, 0, 32*1024),
}

var transport = &Transport{http.DefaultTransport}

func Serve() {
	for _, rp := range config.Config.List {
		target, err := url.Parse(rp.URL)
		if err != nil {
			log.Fatal(err)
		}

		proxy := NewSingleHostReverseProxy(target)
		proxy.Transport = transport

		go http.ListenAndServe(rp.Addr, proxy)
	}
}

func Replay(id uuid.UUID) {
	ts := TsHub.Transactions[id]

	body := ioutil.NopCloser(bytes.NewReader(ts.Req.Body))
	req, err := http.NewRequest(ts.Req.Method, ts.Req.URL.String(), body)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header = CopyHeader(ts.Req.Header)
	_, err = transport.RoundTrip(req)
	if err != nil {
		log.Println(err)
	}
}
