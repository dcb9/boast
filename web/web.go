package web

import (
	"html/template"
	"log"
	"net/http"

	"fmt"
	"github.com/dcb9/boast/config"
	"github.com/dcb9/boast/transaction"
	"github.com/dcb9/boast/web/ws"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/google/uuid"
	"strings"
	"io"
	"bytes"
)

var wsHub = ws.NewHub()
var tsHub = transaction.TsHub

func Serve() {
	go wsHub.Run()
	http.Handle("/static/", http.StripPrefix(
		"/static/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "assets/static"}),
	))

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		tpl, err := Asset("assets/index.html")
		if err != nil {
			log.Fatal(err)
		}

		tmpl := template.Must(
			template.New("index.html").
				Delims("{%", "%}").
				Parse(string(tpl)),
		)
		data := req.Host
		err = tmpl.Execute(rw, data)
		if err != nil {
			log.Fatal(err)
		}
	})
	http.HandleFunc("/responses/", func(rw http.ResponseWriter, req *http.Request) {
		lastSlash := strings.LastIndex(req.URL.Path, "/")
		bs := []byte(req.URL.Path)
		uuidS := bs[lastSlash+1:]

		id, err := uuid.Parse(string(uuidS))
		if err != nil {
			fmt.Fprint(rw, "Bad Request")
		}

		resp := tsHub.Transactions[id].Resp
		src := bytes.NewReader(resp.Body)
		io.Copy(rw, src)
	})

	http.HandleFunc("/ws", func(rw http.ResponseWriter, req *http.Request) {
		ws.Serve(wsHub, rw, req)
	})

	log.Fatal(http.ListenAndServe(config.Config.DebugAddr, nil))
}
