package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"

	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pq" // use wrapped postgres driver

	"github.com/wit-id/blueprint-backend-go/toolkit/db"
)

// NewPostgresDatabase - create & validate postgres connection given certain db.Option
// the caller have the responsibility to close the *sqlx.DB when succeed.
func NewPostgresDatabase(opt *db.Option) (*sql.DB, error) {
	var connStr string // Variabel untuk menyimpan string koneksi akhir

	if opt.ConnectionURL != "" {
		connStr = opt.ConnectionURL
		log.Println("Using DATABASE_URL from db.Option.ConnectionURL")
	} else {
		connURL := &url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(opt.Username, opt.Password),
			Host:   fmt.Sprintf("%s:%d", opt.Host, opt.Port),
			Path:   opt.DatabaseName,
		}
		q := connURL.Query()

		q.Add("sslmode", "require")
		connURL.RawQuery = q.Encode()
		connStr = connURL.String()
		log.Println("Constructing database URL from individual parameters (config.yaml)")
	}

	db, err := apmsql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "postgres: failed to open connection")
	}

	db.SetMaxIdleConns(opt.ConnectionOption.MaxIdle)
	db.SetConnMaxLifetime(opt.ConnectionOption.MaxLifetime)
	db.SetMaxOpenConns(opt.ConnectionOption.MaxOpen)

	ctx, cancel := context.WithTimeout(context.Background(), opt.ConnectionOption.ConnectTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "postgres: failed to ping database: %w")
	}

	log.Println("successfully connected to postgres", connStr)

	go doKeepAliveConnection(db, opt.DatabaseName, opt.ConnectionOption.KeepAliveCheckInterval)

	return db, nil
}

func NewFakePostgresDB() (*sql.DB, error) {
	db, _, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	
	return db, nil
}

func doKeepAliveConnection(db *sql.DB, dbName string, interval time.Duration) {
	for {
		rows, err := db.Query("SELECT 1")
		if err != nil {
			log.Printf("ERROR db.doKeepAliveConnection conn=postgres error=%s db_name=%s\n", err, dbName)
			return
		}

		if rows.Err() != nil {
			log.Printf("ERROR db.doKeepAliveConnection conn=postgres error=%s db_name=%s\n", rows.Err(), dbName)
			return
		}

		if rows.Next() {
			var i int

			_ = rows.Scan(&i)
			log.Printf("SUCCESS db.doKeepAliveConnection counter=%d db_name=%s stats=%+v\n", i, dbName, db.Stats())
		}

		_ = rows.Close()

		time.Sleep(interval)
	}
}
