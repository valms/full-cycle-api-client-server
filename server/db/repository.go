package db

import (
	"context"
	"database/sql"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/valms/full-cycle-api-client-server/server/model"
	"log"
	"time"
)

// SaveExchange Requisito 02 :Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida,
// sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para
// conseguir persistir os dados no banco deverá ser de 10ms.
func SaveExchange(exchange *model.Exchange) error {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancelFunc()

	db, err := sql.Open("sqlite3", "./exchanges.db")

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS exchanges (
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		varBid TEXT,
		pctChange TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		create_date TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.PrepareContext(timeout, "INSERT INTO exchanges "+
		"(code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	for _, currency := range *exchange {
		_, err = stmt.ExecContext(timeout,
			currency.Code,
			currency.Codein,
			currency.Name,
			currency.High,
			currency.Low,
			currency.VarBid,
			currency.PctChange,
			currency.Bid,
			currency.Ask,
			currency.Timestamp,
			currency.CreateDate,
		)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
