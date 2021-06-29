package api

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	sdk_types "github.com/algorand/go-algorand-sdk/types"
	generated "github.com/algorand/indexer/api/generated/v2"
	graphGenerated "github.com/algorand/indexer/api/graph/generated"
	"github.com/algorand/indexer/api/graph/helper"
	"github.com/algorand/indexer/api/graph/model"
	"github.com/algorand/indexer/idb"
)

func (r *accountResolver) AppLocalState(ctx context.Context, obj *model.Account, id uint64) (*model.ApplicationLocalState, error) {
	for _, localState := range obj.AppsLocalState {
		if localState.ID == id {
			return &localState, nil
		}
	}
	return nil, nil
}

func (r *accountResolver) Asset(ctx context.Context, obj *model.Account, id uint64) (*model.AssetHolding, error) {
	for _, holding := range obj.Assets {
		if holding.ID == id {
			return &holding, nil
		}
	}
	return nil, nil
}

func (r *accountResolver) AuthAccount(ctx context.Context, obj *model.Account) (*model.Account, error) {
	if obj.AuthAddr == nil {
		return nil, nil
	}

	return getAccount(r.si, ctx, *obj.AuthAddr)
}

func (r *accountResolver) CreatedApp(ctx context.Context, obj *model.Account, id uint64) (*model.Application, error) {
	for _, app := range obj.CreatedApps {
		if app.ID == id {
			return &app, nil
		}
	}
	return nil, nil
}

func (r *accountResolver) CreatedAsset(ctx context.Context, obj *model.Account, id uint64) (*model.Asset, error) {
	for _, asset := range obj.CreatedAssets {
		if asset.ID == id {
			return &asset, nil
		}
	}
	return nil, nil
}

func (r *applicationLocalStateResolver) Application(ctx context.Context, obj *model.ApplicationLocalState) (*model.Application, error) {
	p := &generated.SearchForApplicationsParams{
		ApplicationId: &obj.ID,
		IncludeAll:    boolPtr(true),
	}
	results, _ := r.si.db.Applications(ctx, p)
	for result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
		return helper.InternalApplicationToModel(result.Application), nil
	}
	return nil, fmt.Errorf("%s: %d", errNoApplicationsFound, obj.ID)
}

func (r *applicationParamsResolver) CreatorAccount(ctx context.Context, obj *model.ApplicationParams) (*model.Account, error) {
	return getAccount(r.si, ctx, obj.Creator)
}

func (r *assetHoldingResolver) Asset(ctx context.Context, obj *model.AssetHolding) (*model.Asset, error) {
	search := generated.SearchForAssetsParams{
		AssetId:    uint64Ptr(obj.ID),
		Limit:      uint64Ptr(1),
		IncludeAll: boolPtr(true),
	}
	options, err := assetParamsToAssetQuery(search)
	if err != nil {
		return nil, err
	}

	assets, _, err := r.si.fetchAssets(ctx, options)
	if err != nil {
		return nil, err
	}

	if len(assets) == 0 {
		return nil, fmt.Errorf("%s: %d", errNoAssetsFound, obj.ID)
	}

	if len(assets) > 1 {
		return nil, fmt.Errorf("%s: %d", errMultipleAssets, obj.ID)
	}

	return helper.InternalAssetToModel(assets[0]), nil
}

func (r *assetParamsResolver) ClawbackAccount(ctx context.Context, obj *model.AssetParams) (*model.Account, error) {
	if obj.Clawback == nil {
		return nil, nil
	}
	return getAccount(r.si, ctx, *obj.Clawback)
}

func (r *assetParamsResolver) CreatorAccount(ctx context.Context, obj *model.AssetParams) (*model.Account, error) {
	return getAccount(r.si, ctx, obj.Creator)
}

func (r *assetParamsResolver) FreezeAccount(ctx context.Context, obj *model.AssetParams) (*model.Account, error) {
	if obj.Freeze == nil {
		return nil, nil
	}
	return getAccount(r.si, ctx, *obj.Freeze)
}

func (r *assetParamsResolver) ManagerAccount(ctx context.Context, obj *model.AssetParams) (*model.Account, error) {
	if obj.Manager == nil {
		return nil, nil
	}
	return getAccount(r.si, ctx, *obj.Manager)
}

func (r *assetParamsResolver) ReserveAccount(ctx context.Context, obj *model.AssetParams) (*model.Account, error) {
	if obj.Reserve == nil {
		return nil, nil
	}
	return getAccount(r.si, ctx, *obj.Reserve)
}

func (r *miniAssetHoldingResolver) Account(ctx context.Context, obj *model.MiniAssetHolding) (*model.Account, error) {
	return getAccount(r.si, ctx, obj.Address)
}

func (r *queryResolver) Block(ctx context.Context, roundNumber uint64) (*model.Block, error) {
	blk, err := r.si.fetchBlock(ctx, roundNumber)
	if err != nil {
		return nil, err
	}

	return helper.InternalBlockToModel(blk), nil
}

func (r *queryResolver) HealthCheck(ctx context.Context) (*model.HealthCheck, error) {
	var errors []string

	health, err := r.si.db.Health()
	if err != nil {
		return nil, fmt.Errorf("problem fetching health: %v", err)
	}

	if health.Error != "" {
		errors = append(errors, fmt.Sprintf("database error: %s", health.Error))
	}

	if r.si.fetcher != nil && r.si.fetcher.Error() != "" {
		errors = append(errors, fmt.Sprintf("fetcher error: %s", r.si.fetcher.Error()))
	}

	result := model.HealthCheck{
		Round:       health.Round,
		IsMigrating: health.IsMigrating,
		DbAvailable: health.DBAvailable,
		Message:     strconv.FormatUint(health.Round, 10),
		Errors:      errors,
	}

	if health.Data != nil {
		result.Data = *health.Data
	}

	return &result, nil
}

func (r *queryResolver) Account(ctx context.Context, address string, includeAll *bool, round *uint64) (*model.AccountResponse, error) {
	addr, errors := decodeAddress(&address, "account-id", make([]string, 0))
	if len(errors) != 0 {
		return nil, fmt.Errorf(errors[0])
	}

	// Special accounts non handling
	isSpecialAccount, err := r.si.isSpecialAccount(address)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", errFailedLoadSpecialAccounts, err)
	}

	if isSpecialAccount {
		return nil, fmt.Errorf(errSpecialAccounts)
	}

	options := idb.AccountQueryOptions{
		EqualToAddress:       addr[:],
		IncludeAssetHoldings: true,
		IncludeAssetParams:   true,
		Limit:                1,
		IncludeDeleted:       boolOrDefault(includeAll),
	}

	accounts, currentRound, err := r.si.fetchAccounts(ctx, options, round)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", errFailedSearchingAccount, err)
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("%s: %s", errNoAccountsFound, address)
	}

	if len(accounts) > 1 {
		return nil, fmt.Errorf("%s: %s", errMultipleAccounts, address)
	}

	return &model.AccountResponse{
		CurrentRound: currentRound,
		Account:      helper.InternalAccountToModel(accounts[0]),
	}, nil
}

func (r *queryResolver) AccountTransactions(ctx context.Context, account string, afterTime *time.Time, assetID *uint64, beforeTime *time.Time, currencyGreaterThan *uint64, currencyLessThan *uint64, limit *uint64, maxRound *uint64, minRound *uint64, next *string, notePrefix []byte, rekeyTo *bool, round *uint64, sigType *model.SigType, txType *model.TxType, id *string) (*model.TransactionsResponse, error) {
	// Check that a valid account was provided
	_, errors := decodeAddress(strPtr(account), "account", make([]string, 0))
	if len(errors) != 0 {
		return nil, fmt.Errorf(errors[0])
	}

	return r.Transactions(ctx, strPtr(account), nil, afterTime, nil, nil, beforeTime, currencyGreaterThan, currencyLessThan, nil, limit, maxRound, minRound, next, notePrefix, rekeyTo, round, sigType, txType, id)
}

func (r *queryResolver) Accounts(ctx context.Context, applicationID *uint64, assetID *uint64, authAddr *string, currencyGreaterThan *uint64, currencyLessThan *uint64, includeAll *bool, limit *uint64, next *string, round *uint64) (*model.AccountsResponse, error) {
	if !r.si.EnableAddressSearchRoundRewind && round != nil {
		return nil, fmt.Errorf(errMultiAcctRewind)
	}

	spendingAddr, errors := decodeAddress(authAddr, "account-id", make([]string, 0))
	if len(errors) != 0 {
		return nil, fmt.Errorf(errors[0])
	}

	options := idb.AccountQueryOptions{
		IncludeAssetHoldings: true,
		IncludeAssetParams:   true,
		Limit:                min(uintOrDefaultValue(limit, defaultAccountsLimit), maxAccountsLimit),
		HasAssetID:           uintOrDefault(assetID),
		HasAppID:             uintOrDefault(applicationID),
		EqualToAuthAddr:      spendingAddr[:],
		IncludeDeleted:       boolOrDefault(includeAll),
	}

	// Set GT/LT on Algos or Asset depending on whether or not an assetID was specified
	if options.HasAssetID == 0 {
		options.AlgosGreaterThan = currencyGreaterThan
		options.AlgosLessThan = currencyLessThan
	} else {
		options.AssetGT = currencyGreaterThan
		options.AssetLT = currencyLessThan
	}

	if next != nil {
		addr, err := sdk_types.DecodeAddress(*next)
		if err != nil {
			return nil, fmt.Errorf(errUnableToParseNext)
		}
		options.GreaterThanAddress = addr[:]
	}

	accounts, currentRound, err := r.si.fetchAccounts(ctx, options, round)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", errFailedSearchingAccount, err)
	}

	var nextToken *string
	if len(accounts) > 0 {
		nextToken = strPtr(accounts[len(accounts)-1].Address)
	}

	return &model.AccountsResponse{
		CurrentRound: currentRound,
		NextToken:    nextToken,
		Accounts:     helper.InternalAccountsToModel(accounts),
	}, nil
}

func (r *queryResolver) Application(ctx context.Context, id uint64, includeAll *bool) (*model.ApplicationResponse, error) {
	p := &generated.SearchForApplicationsParams{
		ApplicationId: &id,
		IncludeAll:    includeAll,
	}
	results, currentRound := r.si.db.Applications(ctx, p)
	out := model.ApplicationResponse{
		CurrentRound: currentRound,
	}
	for result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
		out.Application = helper.InternalApplicationToModel(result.Application)
		return &out, nil
	}
	return nil, fmt.Errorf("%s: %d", errNoApplicationsFound, id)
}

func (r *queryResolver) Applications(ctx context.Context, id *uint64, includeAll *bool, limit *uint64, next *string) (*model.ApplicationsResponse, error) {
	p := &generated.SearchForApplicationsParams{
		ApplicationId: id,
		IncludeAll:    includeAll,
		Limit:         limit,
		Next:          next,
	}
	results, round := r.si.db.Applications(ctx, p)
	apps := make([]generated.Application, 0)
	for result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
		apps = append(apps, result.Application)
	}

	var nextToken *string
	if len(apps) > 0 {
		nextToken = strPtr(strconv.FormatUint(apps[len(apps)-1].Id, 10))
	}

	return &model.ApplicationsResponse{
		Applications: helper.InternalApplicationsToModel(&apps),
		CurrentRound: round,
		NextToken:    nextToken,
	}, nil
}

func (r *queryResolver) Asset(ctx context.Context, id uint64, includeAll *bool) (*model.AssetResponse, error) {
	search := generated.SearchForAssetsParams{
		AssetId:    uint64Ptr(id),
		Limit:      uint64Ptr(1),
		IncludeAll: includeAll,
	}
	options, err := assetParamsToAssetQuery(search)
	if err != nil {
		return nil, err
	}

	assets, currentRound, err := r.si.fetchAssets(ctx, options)
	if err != nil {
		return nil, err
	}

	if len(assets) == 0 {
		return nil, fmt.Errorf("%s: %d", errNoAssetsFound, id)
	}

	if len(assets) > 1 {
		return nil, fmt.Errorf("%s: %d", errMultipleAssets, id)
	}

	return &model.AssetResponse{
		Asset:        helper.InternalAssetToModel(assets[0]),
		CurrentRound: currentRound,
	}, nil
}

func (r *queryResolver) AssetBalances(ctx context.Context, assetID uint64, currencyGreaterThan *uint64, currencyLessThan *uint64, includeAll *bool, limit *uint64, next *string, round *uint64) (*model.AssetBalancesResponse, error) {
	query := idb.AssetBalanceQuery{
		AssetID:        assetID,
		AmountGT:       currencyGreaterThan,
		AmountLT:       currencyLessThan,
		IncludeDeleted: boolOrDefault(includeAll),
		Limit:          min(uintOrDefaultValue(limit, defaultBalancesLimit), maxBalancesLimit),
	}

	if next != nil {
		addr, err := sdk_types.DecodeAddress(*next)
		if err != nil {
			return nil, fmt.Errorf(errUnableToParseNext)
		}
		query.PrevAddress = addr[:]
	}

	balances, currentRound, err := r.si.fetchAssetBalances(ctx, query)
	if err != nil {
		return nil, err
	}

	var nextToken *string
	if len(balances) > 0 {
		nextToken = strPtr(balances[len(balances)-1].Address)
	}

	return &model.AssetBalancesResponse{
		CurrentRound: currentRound,
		NextToken:    nextToken,
		Balances:     helper.InternalMiniAssetHoldingsToModel(&balances),
	}, nil
}

func (r *queryResolver) AssetTransactions(ctx context.Context, address *string, addressRole *model.AddressRole, afterTime *time.Time, assetID uint64, beforeTime *time.Time, currencyGreaterThan *uint64, currencyLessThan *uint64, excludeCloseTo *bool, limit *uint64, maxRound *uint64, minRound *uint64, next *string, notePrefix []byte, rekeyTo *bool, round *uint64, sigType *model.SigType, txType *model.TxType, id *string) (*model.TransactionsResponse, error) {
	return r.Transactions(ctx, address, addressRole, afterTime, nil, uint64Ptr(assetID), beforeTime, currencyGreaterThan, currencyGreaterThan, excludeCloseTo, limit, maxRound, minRound, next, notePrefix, rekeyTo, round, sigType, txType, id)
}

func (r *queryResolver) Assets(ctx context.Context, id *uint64, creator *string, includeAll *bool, limit *uint64, name *string, next *string, unit *string) (*model.AssetsResponse, error) {
	search := generated.SearchForAssetsParams{
		AssetId:    id,
		Limit:      limit,
		IncludeAll: includeAll,
		Next:       next,
		Creator:    creator,
		Name:       name,
		Unit:       unit,
	}
	options, err := assetParamsToAssetQuery(search)
	if err != nil {
		return nil, err
	}

	assets, currentRound, err := r.si.fetchAssets(ctx, options)
	if err != nil {
		return nil, err
	}

	var nextToken *string
	if len(assets) > 0 {
		nextToken = strPtr(strconv.FormatUint(assets[len(assets)-1].Index, 10))
	}

	return &model.AssetsResponse{
		Assets:       helper.InternalAssetsToModel(&assets),
		CurrentRound: currentRound,
		NextToken:    nextToken,
	}, nil
}

func (r *queryResolver) Transaction(ctx context.Context, id string) (*model.TransactionResponse, error) {
	filter, err := transactionParamsToTransactionFilter(generated.SearchForTransactionsParams{
		Txid: strPtr(id),
	})
	if err != nil {
		return nil, err
	}

	// Fetch the transactions
	txns, _, round, err := r.si.fetchTransactions(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", errTransactionSearch, err)
	}

	if len(txns) == 0 {
		return nil, fmt.Errorf("%s: %s", errNoTransactionFound, id)
	}

	if len(txns) > 1 {
		return nil, fmt.Errorf("%s: %s", errMultipleTransactions, id)
	}

	return &model.TransactionResponse{
		CurrentRound: round,
		Transaction:  helper.InternalTransactionToModel(txns[0]),
	}, nil
}

func (r *queryResolver) Transactions(ctx context.Context, address *string, addressRole *model.AddressRole, afterTime *time.Time, applicationID *uint64, assetID *uint64, beforeTime *time.Time, currencyGreaterThan *uint64, currencyLessThan *uint64, excludeCloseTo *bool, limit *uint64, maxRound *uint64, minRound *uint64, next *string, notePrefix []byte, rekeyTo *bool, round *uint64, sigType *model.SigType, txType *model.TxType, id *string) (*model.TransactionsResponse, error) {
	var notePrefixStr *string
	if notePrefix != nil {
		b64 := base64.StdEncoding.EncodeToString(notePrefix)
		notePrefixStr = &b64
	}
	p := generated.SearchForTransactionsParams{
		Limit:               limit,
		Next:                next,
		NotePrefix:          notePrefixStr,
		TxType:              helper.ModalTxTypeToInternal(txType),
		SigType:             helper.ModalSigTypeToInternal(sigType),
		Txid:                id,
		Round:               round,
		MinRound:            minRound,
		MaxRound:            maxRound,
		AssetId:             assetID,
		BeforeTime:          beforeTime,
		AfterTime:           afterTime,
		CurrencyGreaterThan: currencyGreaterThan,
		CurrencyLessThan:    currencyLessThan,
		Address:             address,
		AddressRole:         helper.ModalAddressRoleToInternal(addressRole),
		ExcludeCloseTo:      excludeCloseTo,
		RekeyTo:             rekeyTo,
		ApplicationId:       applicationID,
	}
	filter, err := transactionParamsToTransactionFilter(p)
	if err != nil {
		return nil, err
	}

	// Fetch the transactions
	txns, nextToken, currentRound, err := r.si.fetchTransactions(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", errTransactionSearch, err)
	}

	return &model.TransactionsResponse{
		CurrentRound: currentRound,
		NextToken:    strPtr(nextToken),
		Transactions: helper.InternalTransactionsToModel(&txns),
	}, nil
}

func (r *subscriptionResolver) NewBlock(ctx context.Context) (<-chan *model.Block, error) {
	return r.addBlockListener(ctx), nil
}

func (r *subscriptionResolver) AccountUpdate(ctx context.Context, address string) (<-chan *model.AccountUpdateResponse, error) {
	return r.addAccountListener(ctx, address), nil
}

// Account returns graphGenerated.AccountResolver implementation.
func (r *Resolver) Account() graphGenerated.AccountResolver { return &accountResolver{r} }

// ApplicationLocalState returns graphGenerated.ApplicationLocalStateResolver implementation.
func (r *Resolver) ApplicationLocalState() graphGenerated.ApplicationLocalStateResolver {
	return &applicationLocalStateResolver{r}
}

// ApplicationParams returns graphGenerated.ApplicationParamsResolver implementation.
func (r *Resolver) ApplicationParams() graphGenerated.ApplicationParamsResolver {
	return &applicationParamsResolver{r}
}

// AssetHolding returns graphGenerated.AssetHoldingResolver implementation.
func (r *Resolver) AssetHolding() graphGenerated.AssetHoldingResolver {
	return &assetHoldingResolver{r}
}

// AssetParams returns graphGenerated.AssetParamsResolver implementation.
func (r *Resolver) AssetParams() graphGenerated.AssetParamsResolver { return &assetParamsResolver{r} }

// MiniAssetHolding returns graphGenerated.MiniAssetHoldingResolver implementation.
func (r *Resolver) MiniAssetHolding() graphGenerated.MiniAssetHoldingResolver {
	return &miniAssetHoldingResolver{r}
}

// Query returns graphGenerated.QueryResolver implementation.
func (r *Resolver) Query() graphGenerated.QueryResolver { return &queryResolver{r} }

// Subscription returns graphGenerated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graphGenerated.SubscriptionResolver {
	return &subscriptionResolver{r}
}

type accountResolver struct{ *Resolver }
type applicationLocalStateResolver struct{ *Resolver }
type applicationParamsResolver struct{ *Resolver }
type assetHoldingResolver struct{ *Resolver }
type assetParamsResolver struct{ *Resolver }
type miniAssetHoldingResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
