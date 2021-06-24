package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/algorand/indexer/api/graph/generated"
	"github.com/algorand/indexer/api/graph/model"
)

func (r *queryResolver) Block(ctx context.Context, roundNumber uint64) (*model.Block, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) HealthCheck(ctx context.Context) (*model.HealthCheck, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Account(ctx context.Context, accountID string, includeAll *bool, round *uint64) (*model.AccountResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AccountTransactions(ctx context.Context, accountID string, afterTime *string, assetID *uint64, beforeTime *string, currencyGreaterThan *uint64, currencyLessThan *uint64, limit *uint64, maxRound *uint64, minRound *uint64, next *string, notePrefix *string, rekeyTo *bool, round *uint64, sigType *model.SigType, txType *model.TxType, txid *string) (*model.AccountTransactionsResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Accounts(ctx context.Context, applicationID *uint64, assetID *uint64, authAddr *string, currencyGreaterThan *uint64, currencyLessThan *uint64, includeAll *bool, limit *uint64, next *string, round *uint64) (*model.AccountsResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Application(ctx context.Context, applicationID uint64, includeAll *bool) (*model.ApplicationResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Applications(ctx context.Context, applicationID *uint64, includeAll *bool, limit *uint64, next *string) (*model.ApplicationsResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Asset(ctx context.Context, assetID uint64, includeAll *bool) (*model.AssetResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AssetBalances(ctx context.Context, assetID uint64, currencyGreaterThan *uint64, currencyLessThan *uint64, includeAll *bool, limit *uint64, next *string, round *uint64) (*model.AssetBalancesResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AssetTransactions(ctx context.Context, address *string, addressRole *model.AddressRole, afterTime *string, assetID uint64, beforeTime *string, currencyGreaterThan *uint64, currencyLessThan *uint64, excludeCloseTo *bool, limit *uint64, maxRound *uint64, minRound *uint64, next *string, notePrefix *string, rekeyTo *bool, round *uint64, sigType *model.SigType, txType *model.TxType, txid *string) (*model.AssetTransactionsResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Assets(ctx context.Context, assetID *uint64, creator *string, includeAll *bool, limit *uint64, name *string, next *string, unit *string) (*model.AssetsResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Transaction(ctx context.Context, txid string) (*model.TransactionResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Transactions(ctx context.Context, address *string, addressRole *model.AddressRole, afterTime *string, applicationID *uint64, assetID *uint64, beforeTime *string, currencyGreaterThan *uint64, currencyLessThan *uint64, excludeCloseTo *bool, limit *uint64, maxRound *uint64, minRound *uint64, next *string, notePrefix *string, rekeyTo *bool, round *uint64, sigType *model.SigType, txType *model.TxType, txid *string) (*model.TransactionsResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
