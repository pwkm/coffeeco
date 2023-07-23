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
	Ping(ctx context.Context) error
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
	mongoP := toMongoPurchase(purchase)
	_, err := mr.purchases.InsertOne(ctx, mongoP)
	if err != nil {
		return fmt.Errorf("failed to persis purchase: %w", err)
	}
	return nil
}

type mongoPurchase struct {
	ID                uuid.UUID          `bson:"ID"`
	Store             store.Store        `bson:"Store"`
	ProductToPurchase []coffeeco.Product `bson: "products_purchased"`
	Total             money.Money        `bson: "purchase_total"`
	PaymentMeans      payment.Means      `bson:"payment_means"`
	TimeofPurchase    time.Time          `bson: "created_at"`
	CardToken         *string            `bson: "card_token"`
}

func toMongoPurchase(p Purchase) mongoPurchase {
	return mongoPurchase{
		ID:                p.id,
		Store:             p.Store,
		ProductToPurchase: p.ProductsToPurchase,
		Total:             p.total.Amount(),
		PaymentMeans:      p.PaymentMeans,
		TimeofPurchase:    p.timeofPurchase,
		CardToken:         p.CardToken,
	}
}

func (m mongoPurchase) ToPurchase() Purchase {
	return Purchase{
		id:                 m.ID,
		Store:              m.Store,
		ProductsToPurchase: m.ProductsToPurchase,
		total:              *money.New(m.Total, "USD"),
		PaymentMeans:       m.PaymentMeans,
		timeOfPurchase:     m.TimeOfPurchase,
		CardToken:          m.CardToken,
	}
}

func (mr *MongoRepository) Ping(ctx context.Context) error {
	if _, err := mr.purchases.EstimatedDocumentCount(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}
	return nil
}
