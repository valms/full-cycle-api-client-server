package controller

import (
	"context"
	"encoding/json"
	"github.com/valms/full-cycle-api-client-server/server/client"
	"github.com/valms/full-cycle-api-client-server/server/db"
	"net/http"
	"time"
)

func ExchangeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancelFunc()

	exchange, err := client.GetExchangeFromApi(timeout, request.URL.Query())

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(responseWriter).Encode(&exchange)

	err = db.SaveExchange(exchange)

	if err != nil {
		return
	}

}
