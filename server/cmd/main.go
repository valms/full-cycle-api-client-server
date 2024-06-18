package main

import (
	"github.com/valms/full-cycle-api-client-server/server/controller"
	"net/http"
)

func main() {
	http.HandleFunc("/cotacao", controller.ExchangeHandler)
	http.ListenAndServe(":8080", nil)
}
