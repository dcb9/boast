package boast

import (
	_ "github.com/dcb9/boast/inits/log"
	"net/http/httptest"

	"github.com/dcb9/boast/config"
	"github.com/dcb9/boast/transaction"
	"github.com/dcb9/boast/web"
)

func Serve(s *httptest.Server, addr, debugAddr string) {
	config.Init(s, addr, debugAddr)
	transaction.Serve()
	web.Serve()
}
