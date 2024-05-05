package repository

import (
	"context"
	"errors"

	"github.com/cdleo/api-go-financial-accounting/internal/entity"
	"github.com/cdleo/api-go-financial-accounting/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type accountRepository struct {
	dbClient   *database.MongoDBClient
	collection *mongo.Collection
}

type AccountDTO struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`      /* ObjectID */
	ParentID primitive.ObjectID `bson:"parentId,omitempty"` /* BudgetID or nil */
	entity.AccountHeader
	Balance   float32             `bson:"balance"`   /* Account balance (calculated from movements) */
	Movements []entity.TrxDetails `bson:"movements"` /* Actual transactions in the account */
}

func NewAccountRepository(dbClient *database.MongoDBClient) entity.AccountRepository {
	return &accountRepository{
		dbClient:   dbClient,
		collection: dbClient.Database().Collection("accounts"),
	}
}

func (r *accountRepository) GetByID(ctx context.Context, id string) (*entity.Account, error) {

	var record AccountDTO

	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&record)
	if err != nil {
		return nil, err
	}

	return mapToAccount(&record), nil
}

func (r *accountRepository) GetByHeader(ctx context.Context, header entity.AccountHeader, parentId string) (*entity.Account, error) {

	var record AccountDTO

	filter := bson.D{{"accountheader.type", header.Type}, {"accountheader.name", header.Name}, {"accountheader.currency", header.Currency}}

	// convert parentId string to ObjectId
	if len(parentId) > 0 {
		parentObjId, err := primitive.ObjectIDFromHex(parentId)
		if err != nil {
			return nil, errors.New("invalid id")
		}
		filter = append(filter, bson.E{"parentId", parentObjId})
	}

	err := r.collection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		return nil, err
	}

	return mapToAccount(&record), nil
}

func (r *accountRepository) GetInfo(ctx context.Context, budgetRepo entity.BudgetRepository) ([]*entity.AccountInfo, error) {

	cursor, err := r.collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	results := make([]*entity.AccountInfo, 0)
	for cursor.Next(ctx) {
		var resultDTO AccountDTO
		if err := cursor.Decode(&resultDTO); err != nil {
			return nil, err
		}
		results = append(results, mapToAccountInfo(&resultDTO, budgetRepo))
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *accountRepository) Save(ctx context.Context, value *entity.Account) error {

	result, err := r.collection.InsertOne(ctx, mapToAccountDTO(value))
	if err != nil {
		return err
	}

	value.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *accountRepository) Update(ctx context.Context, value entity.Account) error {

	dto := mapToAccountDTO(&value)

	filter := bson.D{{"_id", dto.ID}}
	result, err := r.collection.ReplaceOne(ctx, filter, dto)
	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return errors.New("not found")
	}
	return nil
}

func (r *accountRepository) UpdateMany(ctx context.Context, values []entity.Account) error {

	sctx, err := r.dbClient.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			r.dbClient.Rollback(sctx)
		}
	}()

	for _, account := range values {
		if err := r.Update(sctx, account); err != nil {
			return err
		}
	}

	return r.dbClient.Commit(sctx)
}
