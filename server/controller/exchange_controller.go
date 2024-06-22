package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/valms/full-cycle-api-client-server/server/client"
	"github.com/valms/full-cycle-api-client-server/server/db"
)

// ExchangeHandler é o manipulador HTTP para a rota /cotacao.
// Ele lida com a requisição, obtém a cotação da API externa, salva no banco de dados e retorna a resposta ao cliente.
func ExchangeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Requisição recebida")

	// Cria um contexto com timeout de 200ms para a chamada da API
	timeout, cancelFunc := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancelFunc()

	// Obtém a cotação da API externa
	exchange, err := client.GetExchangeFromApi(timeout, request.URL.Query())
	if err != nil {
		log.Println("Erro ao obter cotação da API:", err)
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Define o cabeçalho da resposta como JSON
	responseWriter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(responseWriter).Encode(&exchange); err != nil {
		log.Println("Erro ao codificar resposta JSON:", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// Salva a cotação no banco de dados SQLite
	if err := db.SaveExchange(exchange); err != nil {
		log.Println("Erro ao salvar no SQLite:", err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}
