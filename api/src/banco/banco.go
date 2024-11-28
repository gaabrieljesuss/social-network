package banco

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
)

func Conectar() (*sql.DB, error) {
	db, erro := sql.Open("pgx", obterURI())
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil
}

func obterSchema() string {
	return os.Getenv("DATABASE_SCHEMA")
}

func obterURI() string {
	schema := obterSchema()
	user := os.Getenv("POSTGRES_USER")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	name := os.Getenv("POSTGRES_NAME")
	authentication := fmt.Sprintf("%s:%s", user, pwd)
	dst := fmt.Sprintf("%s:%s/%s", host, port, name)
	sslMode := os.Getenv("DATABASE_SSL_MODE")
	return fmt.Sprintf("%s://%s@%s?sslmode=%s", schema, authentication, dst, sslMode)
}
