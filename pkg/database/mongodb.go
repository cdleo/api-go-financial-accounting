package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/cdleo/api-go-financial-accounting/cmd/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MongoDBClient struct {
	dbName  string
	client  *mongo.Client
	options *options.ClientOptions
}

func NewMongoDBClient(cfg config.DBConfig) *MongoDBClient {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	var hosts string
	for _, server := range cfg.Servers {
		if len(hosts) > 0 {
			hosts += ","
		}
		hosts += fmt.Sprintf("%s:%d", server.Host, server.Port)
	}

	var uri string
	if len(cfg.User) > 0 || len(cfg.Password) > 0 {
		uri = fmt.Sprintf("mongodb://%s:%s@%s", cfg.User, cfg.Password, hosts)
	} else {
		uri = fmt.Sprintf("mongodb://%s", hosts)
	}

	if len(cfg.Opts) > 0 {
		uri += fmt.Sprintf("/%s?%s", cfg.DBName, cfg.Opts)
	}

	return &MongoDBClient{
		dbName:  cfg.DBName,
		client:  nil,
		options: options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI),
	}
}

func (c *MongoDBClient) Connect(ctx context.Context) error {
	if c.options == nil {
		return errors.New("client not initialized")
	}

	var err error
	if c.client, err = mongo.Connect(ctx, c.options); err != nil {
		return err
	}

	return c.client.Ping(ctx, readpref.Primary())
}

func (c *MongoDBClient) Disconnect() error {

	var err error
	if c.client != nil {
		err = c.client.Disconnect(context.TODO())
		c.client = nil
	}
	return err
}

func (c *MongoDBClient) Database() *mongo.Database {

	return c.client.Database(c.dbName)
}

func (c *MongoDBClient) StartTransaction(ctx context.Context) (mongo.SessionContext, error) {

	session, err := c.client.StartSession()
	if err != nil {
		return nil, err
	}

	if err = session.StartTransaction(options.Transaction().
		SetReadConcern(readconcern.Snapshot()).
		SetWriteConcern(writeconcern.Majority())); err != nil {
		return nil, err
	}

	return mongo.NewSessionContext(ctx, session), nil
}

func (c *MongoDBClient) Commit(sctx mongo.SessionContext) error {
	defer sctx.EndSession(sctx)
	return sctx.CommitTransaction(sctx)
}

func (c *MongoDBClient) Rollback(sctx mongo.SessionContext) error {
	defer sctx.EndSession(sctx)
	return sctx.AbortTransaction(sctx)
}
