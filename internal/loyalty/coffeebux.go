package loyalty

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/store"

	"github.com/google/uuid"
)

type Coffeebux struct {
	ID                                    uuid.UUID
	store                                 store.Store
	coffeeLover                           coffeeco.CoffeeLover
	FreeDrinksAvailable                   int
	RemainingDrinkPurchasesUntilFreeDrink int
}
