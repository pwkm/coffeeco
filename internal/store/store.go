package store

import (
	coffeeco "coffeeco/internal"

	"github.com/google/uuid"
)

type Store struct {
	ID              uuid.UUID
	Location        string
	ProductsForSale []coffeeco.Product
}
