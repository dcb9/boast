package main

import (
	_ "github.com/dcb9/boast/config"
	"github.com/dcb9/boast/transaction"
	"github.com/dcb9/boast/web"
)

func main() {
	transaction.Serve()
	web.Serve()
}
