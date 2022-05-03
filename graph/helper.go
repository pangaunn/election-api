package graph

import (
	"fmt"
	"math/rand"
)

func generatedKeyForElection(round, IDCard string) string {
	return fmt.Sprintf("e:%s-idcard:%s", round, IDCard)
}
func generatedKeySummaryVote(round, selectedID string) string {
	return fmt.Sprintf("e:%s-candidate:%s-votedcount", round, selectedID)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
