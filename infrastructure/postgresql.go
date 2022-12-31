package infrastructure

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq" // pgx also supported
)

type ConnDB struct {
	pool *sql.DB
}

func Open() (*sql.DB, error) {
	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME") + "?sslmode=disable"
	pool, err := sql.Open("postgres", "postgres://"+dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(); err != nil {
		return nil, err
	}

	return pool, nil
}

func (p *ConnDB) Close() {
	p.pool.Close()
}
