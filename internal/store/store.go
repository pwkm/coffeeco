package store

import "github.com/google/uuid"

type Store struct {
	ID       uuid.UUID
	Location string
}
