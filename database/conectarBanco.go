package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "ep-gentle-night-ap1246mx-pooler.c-7.us-east-1.aws.neon.tech"
	port     = 5432
	user     = "neondb_owner"
	password = "npg_dY4ehmULb5NV"
	dbname   = "neondb"
)

func ConectarBanco() (*sql.DB, error) {
	stringConexao := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", stringConexao)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Conectado ao Banco de Dados: ", dbname)
	return db, err
}
