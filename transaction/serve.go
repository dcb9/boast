package transaction

import (
	"log"
	"net/http"
	"net/url"

	"github.com/dcb9/boast/config"
)

var TxHub *Hub

var transport = &Transport{http.DefaultTransport}

func Serve() {
	TxHub = NewTxHub()
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
