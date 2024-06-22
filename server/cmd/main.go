package main

import (
	"log"
	"net/http"

	"github.com/valms/full-cycle-api-client-server/server/controller"
)

// main é a função principal que configura e inicia o servidor HTTP.
func main() {
	// Define o manipulador para a rota /cotacao
	http.HandleFunc("/cotacao", controller.ExchangeHandler)

	log.Println("Servidor rodando na porta 8080")

	// Inicia o servidor HTTP na porta 8080 e loga qualquer erro que ocorra
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
