package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/valms/full-cycle-api-client-server/server/model"
)

// main é a função principal que orquestra a execução do programa.
func main() {
	// Cria um contexto com timeout de 300ms
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Obter a cotação
	responseData, err := fetchExchangeRate(ctx)
	if err != nil {
		log.Println("Erro ao obter cotação:", err)
		return
	}

	// Pegar a primeira chave do mapa da resposta
	firstMapKey := returnFirstMapKey(responseData)
	log.Println("Primeira chave é:", firstMapKey)

	bid := responseData[firstMapKey].Bid
	log.Println(bid)

	// Salvar em um arquivo
	err = saveExchangeRateToFile(firstMapKey, bid)
	if err != nil {
		log.Println("Erro ao salvar cotação em arquivo:", err)
		return
	}

	log.Println("Cotação salva com sucesso!")
}

// fetchExchangeRate faz uma requisição HTTP para obter a cotação e retorna os dados decodificados.
func fetchExchangeRate(ctx context.Context) (model.Exchange, error) {
	resp, err := doRequest(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Erro: tempo limite excedido ao fazer requisição")
		} else {
			log.Println("Erro ao fazer requisição:", err)
		}
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		bodyErr := Body.Close()
		if bodyErr != nil {
			log.Println("Erro ao fechar o Body:", bodyErr)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Println("Resposta recebida com erro:", resp.Status)
		return nil, fmt.Errorf("resposta recebida com erro: %s", resp.Status)
	}

	var responseData model.Exchange
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		log.Println("Erro ao ler resposta:", err)
		return nil, err
	}

	return responseData, nil
}

// saveExchangeRateToFile salva a cotação em um arquivo de texto.
func saveExchangeRateToFile(key, bid string) error {
	formattedKey := divideKey(key)
	content := fmt.Sprintf("Cotação %s: %s", formattedKey, bid)
	err := os.WriteFile("cotacao.txt", []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

// returnFirstMapKey retorna a primeira chave de um mapa do tipo model.Exchange.
func returnFirstMapKey(m model.Exchange) string {
	for k := range m {
		return k
	}
	return ""
}

// divideKey formata a chave adicionando um hífen entre os primeiros três caracteres e os restantes.
func divideKey(key string) string {
	return key[:3] + "-" + key[3:]
}

// doRequest cria e envia uma requisição HTTP para obter a cotação.
func doRequest(ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Println("Erro ao criar requisição:", err)
		return nil, err
	}

	queryValues := req.URL.Query()
	queryValues.Set("moedaOrigem", "BRL")
	queryValues.Set("moedaDestino", "USD")
	req.URL.RawQuery = queryValues.Encode()

	log.Println(req.URL.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Erro: tempo limite excedido ao fazer requisição")
		} else {
			log.Println("Erro ao fazer requisição:", err)
		}
		return nil, err
	}

	return resp, nil
}
