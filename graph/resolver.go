package graph

import (
	"election-api/graph/model"
	"sync"

	"github.com/go-redis/redis"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MU        sync.Mutex
	Observers map[string]chan *model.CandidateVoteUpdated
	Cache     *redis.Client
}
