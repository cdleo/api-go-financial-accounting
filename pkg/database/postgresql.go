package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cdleo/api-go-financial-accounting/cmd/config"
	"github.com/cdleo/go-commons/logger"
	"github.com/cdleo/go-sqldb"
	"github.com/cdleo/go-sqldb/adapter"
	"github.com/cdleo/go-sqldb/connector"
)

type PostgreSQLClient struct {
	proxy *sqldb.SQLProxy
	db    *sql.DB
}

func NewPostgreSQLClient(cfg config.DBConfig) *PostgreSQLClient {

	proxy := sqldb.NewSQLProxyBuilder(connector.NewPostgreSqlConnector(cfg.Servers[0].Host, cfg.Servers[0].Port, cfg.User, cfg.Password, cfg.DBName)).
		WithAdapter(adapter.NewNoopAdapter()).
		WithLogger(logger.NewNoLogLogger()).
		Build()

	return &PostgreSQLClient{
		proxy: proxy,
		db:    nil,
	}
}

func (c *PostgreSQLClient) Connect(ctx context.Context) error {

	var err error
	c.db, err = c.proxy.Open()
	if err != nil {
		fmt.Printf("Unable to connect to DB. Error: %v", err)
	}
	return err
}

func (c *PostgreSQLClient) Disconnect() error {

	var err error
	if c.db != nil {
		err = c.proxy.Close()
		c.db = nil
	}
	return err
}

func (c *PostgreSQLClient) Database() *sql.DB {
	return c.db
}
