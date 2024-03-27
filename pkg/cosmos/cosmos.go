package cosmos

import (
	"context"

	"github.com/evmos/validator-status/pkg/database"
)

type Cosmos struct {
	apiURL string
	rpcURL string
	db     *database.Queries
	ctx    context.Context
}

func NewCosmos(apiURL string, rpcURL string, db *database.Queries) *Cosmos {
	return &Cosmos{
		apiURL: apiURL,
		rpcURL: rpcURL,
		db:     db,
		ctx:    context.Background(),
	}
}
