package graph

import (
	"context"
	"fmt"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

func ValidIDCard(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	IDCard := fmt.Sprintf("%v", ctx.Value("IDCard"))
	if len(IDCard) != 13 {
		return nil, fmt.Errorf("Access denied")
	}

	if _, err := strconv.Atoi(IDCard); err != nil {
		return nil, fmt.Errorf("Access denied")
	}

	sum := 0
	for i := 0; i < len(IDCard)-1; i++ {
		n, _ := strconv.Atoi(string(IDCard[i]))
		sum = sum + ((13 - i) * n)
	}

	lastNumber, _ := strconv.Atoi(string(IDCard[len(IDCard)-1]))
	checksum := (11 - (sum % 11)) % 10

	if lastNumber != checksum {
		return nil, fmt.Errorf("Access denied")
	}

	return next(ctx)
}
