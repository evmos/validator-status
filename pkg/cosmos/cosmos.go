package cosmos

import (
	"context"

	"github.com/evmos/validator-status/pkg/database"
)

type Cosmos struct {
	apiURL string
	db     *database.Queries
	ctx    context.Context
}

func NewCosmos(apiURL string, db *database.Queries) *Cosmos {
	return &Cosmos{
		apiURL: apiURL,
		db:     db,
		ctx:    context.Background(),
	}
}
