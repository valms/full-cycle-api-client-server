package model

// Exchange representa um mapa de moedas, onde a chave é uma string e o valor é um objeto Currency.
type Exchange map[string]Currency

// Currency representa os detalhes de uma moeda específica.
type Currency struct {
	Code       string `json:"code"`        // Código da moeda (ex: USD)
	Codein     string `json:"codein"`      // Código da moeda de conversão (ex: BRL)
	Name       string `json:"name"`        // Nome da moeda (ex: Dólar Americano)
	High       string `json:"high"`        // Maior valor da moeda no período
	Low        string `json:"low"`         // Menor valor da moeda no período
	VarBid     string `json:"varBid"`      // Variação do valor da moeda
	PctChange  string `json:"pctChange"`   // Percentual de mudança do valor da moeda
	Bid        string `json:"bid"`         // Valor de compra da moeda
	Ask        string `json:"ask"`         // Valor de venda da moeda
	Timestamp  string `json:"timestamp"`   // Timestamp da cotação
	CreateDate string `json:"create_date"` // Data de criação da cotação
}
