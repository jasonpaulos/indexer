package api

import (
	"context"
	"fmt"

	"github.com/algorand/indexer/api/graph/helper"
	"github.com/algorand/indexer/api/graph/model"
	"github.com/algorand/indexer/idb"
)

func getAccount(si *ServerImplementation, ctx context.Context, address string) (*model.Account, error) {
	addr, errors := decodeAddress(&address, "address", make([]string, 0))
	if len(errors) != 0 {
		return nil, fmt.Errorf(errors[0])
	}

	// Special accounts non handling
	isSpecialAccount, err := si.isSpecialAccount(address)
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
		IncludeDeleted:       true,
	}

	accounts, _, err := si.fetchAccounts(ctx, options, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", errFailedSearchingAccount, err)
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("%s: %s", errNoAccountsFound, address)
	}

	if len(accounts) > 1 {
		return nil, fmt.Errorf("%s: %s", errMultipleAccounts, address)
	}

	return helper.InternalAccountToModel(accounts[0]), nil
}
