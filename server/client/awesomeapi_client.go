package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valms/full-cycle-api-client-server/server/model"
	"net/http"
	"net/url"
)

// GetExchangeFromApi Requisito 01: O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
// O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL
// e em seguida deverá retornar no formato JSON o resultado para o cliente.
func GetExchangeFromApi(ctx context.Context, params url.Values) (*model.Exchange, error) {
	if len(params["moedaOrigem"]) == 0 {
		return nil, errors.New("Parâmetro 'moedaOrigem' não encontrado!")
	}

	if len(params["moedaDestino"]) == 0 {
		return nil, errors.New("Parâmetro 'moedaDestino' não encontrado!")
	}

	code := params["moedaOrigem"][0]
	codeIn := params["moedaDestino"][0]

	client := &http.Client{}
	requestWithContext, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://economia.awesomeapi.com.br/json/last/%s-%s", code, codeIn), nil)

	if err != nil {
		return nil, err
	}

	response, err := client.Do(requestWithContext)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var exchange model.Exchange
	err = json.NewDecoder(response.Body).Decode(&exchange)

	if err != nil {
		return nil, err
	}

	return &exchange, nil
}
