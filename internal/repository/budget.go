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

type budgetRepository struct {
	dbClient    *database.MongoDBClient
	collection  *mongo.Collection
	accountRepo entity.AccountRepository
}

type BudgetDTO struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"` /* ObjectID */
	Type        entity.BudgetType   `bson:"type"`          /* Monthly or Travel */
	Month       int                 `bson:"month"`         /* Month of the Budget */
	Year        int                 `bson:"year"`          /* Year of the Budget */
	Description string              `bson:"description"`   /* Optional: Descriptive text*/
	Accounts    []*BudgetAccountDTO `bson:"accounts"`
}

type BudgetAccountDTO struct {
	AccountID primitive.ObjectID       `bson:"account_id,omitempty"` /* ObjectID */
	Expected  []entity.ExpenseExpected `bson:"expected,omitempty"`
}

func NewBudgetRepository(dbClient *database.MongoDBClient, accountRepo entity.AccountRepository) entity.BudgetRepository {
	return &budgetRepository{
		dbClient:    dbClient,
		collection:  dbClient.Database().Collection("budgets"),
		accountRepo: accountRepo,
	}
}

func (r *budgetRepository) GetByDate(ctx context.Context, month int, year int) (*entity.Budget, error) {

	var record *BudgetDTO = nil

	filter := bson.D{{"month", month}, {"year", year}}

	err := r.collection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		return nil, err
	}
	return mapToBudget(record, r.accountRepo)
}

func (r *budgetRepository) GetById(ctx context.Context, id string) (*entity.Budget, error) {

	var record *BudgetDTO = nil

	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&record)
	if err != nil {
		return nil, err
	}
	return mapToBudget(record, r.accountRepo)
}

func (r *budgetRepository) GetAll(ctx context.Context) ([]*entity.Budget, error) {

	cursor, err := r.collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	results := make([]*entity.Budget, 0)
	for cursor.Next(ctx) {
		var resultDTO BudgetDTO
		if err := cursor.Decode(&resultDTO); err != nil {
			return nil, err
		}
		if result, err := mapToBudget(&resultDTO, r.accountRepo); err != nil {
			return nil, err
		} else {
			results = append(results, result)
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *budgetRepository) Save(ctx context.Context, value *entity.Budget) error {

	sctx, err := r.dbClient.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			r.dbClient.Rollback(sctx)
		}
	}()

	value.ID = primitive.NewObjectID().Hex()

	for i := 0; i < len(value.Accounts); i++ {
		value.Accounts[i].Account.ParentID = value.ID
		err = r.accountRepo.Save(sctx, &value.Accounts[i].Account)
		if err != nil {
			return err
		}
	}

	dto := mapToBudgetDTO(value)

	_, err = r.collection.InsertOne(sctx, dto)
	if err != nil {
		return err
	}

	return r.dbClient.Commit(sctx)
}

func (r *budgetRepository) Update(ctx context.Context, value *entity.Budget) error {

	dto := mapToBudgetDTO(value)

	filter := bson.D{{"_id", dto.ID}}
	result, err := r.collection.ReplaceOne(ctx, filter, dto)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("not found")
	}

	return nil
}
