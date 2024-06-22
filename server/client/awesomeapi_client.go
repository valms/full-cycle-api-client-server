package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/valms/full-cycle-api-client-server/server/model"
)

// GetExchangeFromApi realiza uma requisição HTTP para obter a cotação de uma moeda.
// Requisito 01: O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
// O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL
// e em seguida deverá retornar no formato JSON o resultado para o cliente.
func GetExchangeFromApi(ctx context.Context, params url.Values) (*model.Exchange, error) {
	// Verifica os parâmetros da URL
	if err := validateParams(params); err != nil {
		return nil, err
	}

	// Cria a requisição HTTP
	req, err := createRequest(ctx, params)
	if err != nil {
		return nil, err
	}

	// Executa a requisição HTTP
	resp, err := executeRequest(req)
	if err != nil {
		return nil, err
	}
	defer closeResponseBody(resp.Body)

	// Decodifica a resposta JSON
	exchange, err := decodeResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	return exchange, nil
}

// validateParams verifica se os parâmetros 'moedaOrigem' e 'moedaDestino' estão presentes.
func validateParams(params url.Values) error {
	if len(params["moedaOrigem"]) == 0 {
		return errors.New("parâmetro 'moedaOrigem' não encontrado")
	}
	if len(params["moedaDestino"]) == 0 {
		return errors.New("parâmetro 'moedaDestino' não encontrado")
	}
	return nil
}

// createRequest cria uma requisição HTTP com contexto.
func createRequest(ctx context.Context, params url.Values) (*http.Request, error) {
	code := params["moedaOrigem"][0]
	codeIn := params["moedaDestino"][0]
	requestUrl := fmt.Sprintf("https://economia.awesomeapi.com.br/json/last/%s-%s", code, codeIn)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		log.Println("Erro ao criar requisição ao awesomeapi:", err)
		return nil, err
	}
	return req, nil
}

// executeRequest executa a requisição HTTP e retorna a resposta.
func executeRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Erro: tempo limite excedido ao chamar a API de cotação do dólar")
		} else {
			log.Println("Erro ao chamar o awesomeapi:", err)
		}
		return nil, err
	}
	return resp, nil
}

// closeResponseBody fecha o corpo da resposta HTTP.
func closeResponseBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Println("Erro ao fechar o Body:", err)
	}
}

// decodeResponse decodifica a resposta JSON e retorna um objeto Exchange.
func decodeResponse(body io.Reader) (*model.Exchange, error) {
	var exchange model.Exchange
	if err := json.NewDecoder(body).Decode(&exchange); err != nil {
		log.Println("Erro ao decodificar resposta JSON:", err)
		return nil, err
	}
	return &exchange, nil
}
