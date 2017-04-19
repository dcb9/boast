package main

import (
	"github.com/dcb9/boast/config"
	"github.com/dcb9/boast/transaction"
	"github.com/dcb9/boast/web"
)

func main() {
	config.CmdInit()
	transaction.Serve()
	web.Serve()
}
