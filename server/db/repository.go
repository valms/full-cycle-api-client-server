package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/valms/full-cycle-api-client-server/server/model"
)

// SaveExchange salva a cotação no banco de dados SQLite.
// Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida,
// sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para
// conseguir persistir os dados no banco deverá ser de 10ms.
func SaveExchange(exchange *model.Exchange) error {
	// Cria um contexto com timeout de 10ms para a operação de salvamento no banco de dados
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()

	// Abre uma conexão com o banco de dados SQLite
	db, databaseErr := openDatabase()
	if databaseErr != nil {
		return databaseErr
	}
	defer closeDatabase(db)

	// Cria a tabela exchanges se ela não existir
	if tableErr := createTable(db); tableErr != nil {
		return tableErr
	}

	// Prepara o statement para inserção de dados
	stmt, tableErr := prepareStatement(db, timeoutCtx)
	if tableErr != nil {
		return tableErr
	}
	defer closeStatement(stmt)

	// Insere cada cotação no banco de dados
	if insertErr := insertExchange(stmt, timeoutCtx, exchange); insertErr != nil {
		return insertErr
	}

	return nil
}

// openDatabase abre uma conexão com o banco de dados SQLite.
func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./exchanges.db")
	if err != nil {
		log.Println("Erro ao abrir conexão:", err)
		return nil, err
	}
	return db, nil
}

// closeDatabase fecha a conexão com o banco de dados SQLite.
func closeDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Println("Erro ao fechar conexão com o banco de dados:", err)
	}
}

// createTable cria a tabela exchanges se ela não existir.
func createTable(db *sql.DB) error {
	_, err := db.Exec(`
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
		log.Println("Erro ao criar tabela:", err)
		return err
	}
	return nil
}

// prepareStatement prepara o statement para inserção de dados.
func prepareStatement(db *sql.DB, ctx context.Context) (*sql.Stmt, error) {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO exchanges
    (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)  
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Println("Erro ao preparar statement:", err)
		return nil, err
	}
	return stmt, nil
}

// closeStatement fecha o statement.
func closeStatement(stmt *sql.Stmt) {
	if err := stmt.Close(); err != nil {
		log.Println("Erro ao fechar statement:", err)
	}
}

// insertExchange insere cada cotação no banco de dados.
func insertExchange(stmt *sql.Stmt, ctx context.Context, exchange *model.Exchange) error {
	for _, currency := range *exchange {
		_, err := stmt.ExecContext(ctx,
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
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				log.Println("Erro: tempo limite excedido ao salvar cotação no banco de dados")
			} else {
				log.Println("Erro ao executar statement:", err)
			}
			return err
		}
	}
	return nil
}
