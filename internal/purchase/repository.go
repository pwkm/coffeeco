package purchase

import "context"

type repository interface {
	Store(ctx context.Context, purchase Purchase) error
}
