package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"election-api/graph/generated"
	"election-api/graph/model"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func (r *mutationResolver) Vote(ctx context.Context, id string) (bool, error) {
	round := r.Cache.Get("e:current:active")
	err := round.Err()
	if err != nil {
		return false, fmt.Errorf("Not open election at the moment")
	}

	IDCard := fmt.Sprintf("%v", ctx.Value("IDCard"))
	voteKey := generatedKeyForElection(round.Val(), IDCard)
	boolCmd := r.Cache.SetNX(voteKey, 1, time.Minute*60)

	result, err := boolCmd.Result()
	if err != nil {
		return false, fmt.Errorf("Something wrong with server")
	}

	if !result {
		return false, fmt.Errorf("Duplicated IDCard")
	}

	summaryKey := generatedKeySummaryVote(round.Val(), id)
	intCmd := r.Cache.Incr(summaryKey)

	incrResult, err := intCmd.Result()
	if err != nil {
		return false, fmt.Errorf("Something wrong with server")
	}

	if incrResult == 1 {
		r.Cache.Expire(summaryKey, time.Minute*60)
	}

	r.MU.Lock()
	votedCount := int(incrResult)
	for _, observer := range r.Observers {
		observer <- &model.CandidateVoteUpdated{
			ID:         id,
			VotedCount: votedCount,
		}
	}
	r.MU.Unlock()
	return true, nil
}

func (r *mutationResolver) Open(ctx context.Context) (bool, error) {
	stringCmd := r.Cache.Get("e:current:active")
	active, err := stringCmd.Result()
	if err != nil && err.Error() != "redis: nil" {
		return false, err
	}

	if active == "" {
		cmd := r.Cache.Get("e:current")
		roundPrevious, _ := cmd.Result()
		if roundPrevious == "" {
			r.Cache.Set("e:current", 1, -1)
			r.Cache.Set("e:current:active", 1, -1)
		} else {
			incr := r.Cache.Incr("e:current")
			newActive := strconv.Itoa(int(incr.Val()))
			r.Cache.Set("e:current:active", newActive, -1)
		}
	} else {
		return false, fmt.Errorf("invalid open election new round because round: %s activeted ", active)
	}
	return true, nil
}

func (r *mutationResolver) Close(ctx context.Context) (bool, error) {
	cmd := r.Cache.Del("e:current:active")
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return true, nil
}

func (r *queryResolver) Candidates(ctx context.Context) ([]*model.Candidate, error) {
	stringSliceCmd := r.Cache.Keys("candidate:*")
	keys, err := stringSliceCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("Something wrong with server")
	}

	sliceCmd := r.Cache.MGet(keys...)
	slice, err := sliceCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("Something wrong with server")
	}

	var candidates []*model.Candidate

	roundCmd := r.Cache.Get("e:current")
	round, err := roundCmd.Result()
	if err != nil && err.Error() != "redis: nil" {
		return nil, fmt.Errorf("Something wrong with server")
	}

	for _, s := range slice {
		var candidate model.Candidate
		str := fmt.Sprintf("%v", s)
		json.Unmarshal([]byte(str), &candidate)
		candidates = append(candidates, &candidate)

		// Can use Mget to improve perf
		summaryKey := generatedKeySummaryVote(round, candidate.ID)
		stringCmd := r.Cache.Get(summaryKey)
		str, err = stringCmd.Result()
		if err != nil && err.Error() != "redis: nil" {
			return nil, fmt.Errorf("Something wrong with server")
		}
		voteCount, _ := strconv.Atoi(str)
		candidate.VotedCount = voteCount
	}

	return candidates, nil
}

func (r *queryResolver) Candidate(ctx context.Context, id string) (*model.Candidate, error) {
	key := fmt.Sprintf("candidate:%s", id)
	stringCmd := r.Cache.Get(key)
	str, err := stringCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("Candidate %s not found", id)
	}

	var candidate model.Candidate
	json.Unmarshal([]byte(str), &candidate)

	roundCmd := r.Cache.Get("e:current")
	round, _ := roundCmd.Result()

	summaryKey := generatedKeySummaryVote(round, id)
	stringCmd = r.Cache.Get(summaryKey)
	str, err = stringCmd.Result()
	if err != nil && err.Error() != "redis: nil" {
		return nil, fmt.Errorf("Something wrong with server")
	}

	voteCount, _ := strconv.Atoi(str)
	candidate.VotedCount = voteCount

	return &candidate, nil
}

func (r *subscriptionResolver) VoteUpdated(ctx context.Context) (<-chan *model.CandidateVoteUpdated, error) {
	id := randString(8)
	events := make(chan *model.CandidateVoteUpdated, 1)

	go func() {
		<-ctx.Done()
		r.MU.Lock()
		delete(r.Observers, id)
		r.MU.Unlock()
	}()

	r.MU.Lock()
	r.Observers[id] = events
	r.MU.Unlock()

	return events, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
