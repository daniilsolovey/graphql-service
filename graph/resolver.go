package graph

import (
	"github.com/daniilsolovey/graphql-service/internal/operator"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Operator operator.Operator
}
