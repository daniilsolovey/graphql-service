package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/daniilsolovey/graphql-service/graph/generated"
	"github.com/daniilsolovey/graphql-service/graph/model"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

const (
	INTERNAL_ERROR = "Внутренняя ошибка. Попробуйте позже"
)

func (r *mutationResolver) RequestSignInCode(
	ctx context.Context,
	input model.RequestSignInCodeInput,
) (*model.ErrorPayload, error) {
	// in future you can add checking phone number length (=8), start number (+7)
	if input.Phone == "" {
		return &model.ErrorPayload{Message: "phone number required"}, nil
	}

	err := r.Operator.RequestSignInCode(input.Phone)
	if err != nil {
		log.Error(err)
		return &model.ErrorPayload{Message: INTERNAL_ERROR}, nil
	}

	return nil, nil
}

func (r *mutationResolver) SignInByCode(
	ctx context.Context,
	input model.SignInByCodeInput,
) (model.SignInOrErrorPayload, error) {
	result, err := r.Operator.SignInByCode(input.Phone, input.Code)
	if err != nil {
		log.Error(err)
		return &model.ErrorPayload{Message: INTERNAL_ERROR}, nil
	}

	return result, nil
}

func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	result, err := r.Operator.GetAllProducts()
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to get all products",
		)
	}

	return result, nil
}

func (r *queryResolver) Viewer(ctx context.Context) (*model.Viewer, error) {
	token := ctx.Value("Authorization").(string)
	result, err := r.Operator.Viewer(token)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to view user",
		)
	}

	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
