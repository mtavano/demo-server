package main

import (
	"log"
	"net/http"
	"os"

	"context"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// Store is the database wrapper
type SqlStore struct {
	*sqlx.DB
}

func NewSqlStore(driver, dsn string) (*SqlStore, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// Maximum Idle Connections
	db.SetMaxIdleConns(20)
	// Idle Connection Timeout
	db.SetConnMaxIdleTime(1 * time.Second)
	// Connection Lifetime
	db.SetConnMaxLifetime(30 * time.Second)

	_, err = db.Exec("SELECT true;")
	if err != nil {
		return nil, err
	}

	return &SqlStore{db}, nil
}

func (st *SqlStore) BeginTx(ctx context.Context) (any, error) {
	tx, err := st.DB.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, errors.Wrap(err, "database: Store.BeginTx st.BeginTxx error")
	}
	return tx, nil
}

//func (st *SqlStore) Exec(query string, params ...interface{}) (sql.Result, error) {
//return st.Exec(query, params...)
//}

func main() {
	log.Println("Start demo server")
	databaseDSN := os.Getenv("DATABASE_DSN")
	db, err := NewSqlStore("postgres", databaseDSN)
	if err != nil {
		panic(err)
	}
	log.Println("db:", db)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
