// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	GetHappycatFact(ctx context.Context, id uuid.UUID) (*HappycatFact, error)
	ListHappycatFacts(ctx context.Context) ([]*HappycatFact, error)
	SaveHappycatFact(ctx context.Context, arg SaveHappycatFactParams) error
}

var _ Querier = (*Queries)(nil)
