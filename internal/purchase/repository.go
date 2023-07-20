package purchase

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/payment"
	"coffeeco/internal/store"
	"context"
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ------------------------ REPOSITORY INTERFACE --------------------------------------
//
// ------------------------------------------------------------------------------------
type Repository interface {
	Store(ctx context.Context, purchase Purchase) error
}

// --------------------------- MONGO REPOSITORY  --------------------------------------
//
// ------------------------------------------------------------------------------------
type MongoRepository struct {
	purchases *mongo.Collection
}

func NewMongoRepo(ctx context.Context, connectionString string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}
	purchases := client.Database("coffeeco").Collection("purchases")

	return &MongoRepository{
		purchases: purchases,
	}, nil
}

func (mr *MongoRepository) Store(ctx context.Context, purchase Purchase) error {
	mongoP := New(purchase)
	_, err := mr.purchases.InsertOne(ctx, mongoP)
	if err != nil {
		return fmt.Errorf("failed to persis purchase: %w", err)
	}
	return nil
}

type mongoPurchase struct {
	id                uuid.UUID
	store             store.Store
	productToPurchase []coffeeco.Product
	total             money.Money
	PaymentMeans      payment.Means
	timeofPurchase    time.Time
	cardToken         *string
}

func toMongoPurchase(p Purchase) mongoPurchase {
	return mongoPurchase{
		id:                p.id,
		store:             p.Store,
		productToPurchase: p.ProductsToPurchase,
		total:             p.total,
		PaymentMeans:      p.PaymentMeans,
		timeofPurchase:    p.timeofPurchase,
		cardToken:         p.cardToken,
	}
}
