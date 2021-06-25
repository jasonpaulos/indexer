package helper

import (
	"encoding/base64"
	"fmt"

	"github.com/algorand/indexer/api/generated/v2"
	"github.com/algorand/indexer/api/graph/model"
)

func boolPtrOrDefault(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func uint64PtrOrDefault(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}

func byteslicePtrOrDefault(p *[]byte) []byte {
	if p == nil {
		return []byte{}
	}
	return *p
}

func InternalSigTypeToModel(sigType *string) *model.SigType {
	if sigType == nil {
		return nil
	}

	var result model.SigType

	switch *sigType {
	case "sig":
		result = model.SigTypeSig
	case "msig":
		result = model.SigTypeMsig
	case "lsig":
		result = model.SigTypeLsig
	default:
		panic(fmt.Errorf("Unexpected sig type: %s", *sigType))
	}

	return &result
}

func InternalAccountParticipationToModel(participation *generated.AccountParticipation) *model.AccountParticipation {
	if participation == nil {
		return nil
	}

	return &model.AccountParticipation{
		SelectionParticipationKey: participation.SelectionParticipationKey,
		VoteFirstValid:            participation.VoteFirstValid,
		VoteKeyDilution:           participation.VoteKeyDilution,
		VoteLastValid:             participation.VoteLastValid,
		VoteParticipationKey:      participation.VoteParticipationKey,
	}
}

func InternalAccountsToModel(accounts []generated.Account) []model.Account {
	converted := make([]model.Account, len(accounts))
	for i, account := range accounts {
		converted[i] = *InternalAccountToModel(account)
	}
	return converted
}

func InternalAccountToModel(account generated.Account) *model.Account {
	return &model.Account{
		Address:                     account.Address,
		Amount:                      account.Amount,
		AmountWithoutPendingRewards: account.AmountWithoutPendingRewards,
		AppsLocalState:              InternalApplicationLocalStatesToModel(account.AppsLocalState),
		AppsTotalExtraPages:         uint64PtrOrDefault(account.AppsTotalExtraPages),
		AppsTotalSchema:             InternalApplicationStateSchemaToModel(account.AppsTotalSchema),
		Assets:                      InternalAssetHoldingsToModel(account.Assets),
		AuthAddr:                    account.AuthAddr,
		ClosedAtRound:               account.ClosedAtRound,
		CreatedApps:                 InternalApplicationsToModel(account.CreatedApps),
		CreatedAssets:               InternalAssetsToModel(account.CreatedAssets),
		CreatedAtRound:              account.CreatedAtRound,
		Deleted:                     boolPtrOrDefault(account.Deleted),
		Participation:               InternalAccountParticipationToModel(account.Participation),
		PendingRewards:              account.PendingRewards,
		RewardBase:                  account.RewardBase,
		Rewards:                     account.Rewards,
		Round:                       account.Round,
		SigType:                     InternalSigTypeToModel(account.SigType),
		Status:                      account.Status,
	}
}

func InternalApplicationStateSchemaToModel(schema *generated.ApplicationStateSchema) *model.ApplicationStateSchema {
	if schema == nil {
		return nil
	}

	return &model.ApplicationStateSchema{
		NumByteSlice: schema.NumByteSlice,
		NumUint:      schema.NumUint,
	}
}

func InternalTealKeyValueStoreToModel(store *generated.TealKeyValueStore) []model.TealKeyValue {
	if store == nil {
		return []model.TealKeyValue{}
	}

	converted := make([]model.TealKeyValue, len(*store))
	for i, keyValue := range *store {
		converted[i] = InternalTealKeyValueToModel(keyValue)
	}
	return converted
}

func InternalTealKeyValueToModel(keyValue generated.TealKeyValue) model.TealKeyValue {
	keyBytes, err := base64.RawStdEncoding.DecodeString(keyValue.Key)
	if err != nil {
		panic(err)
	}
	return model.TealKeyValue{
		Key:   keyBytes,
		Value: InternalTealValueToModel(keyValue.Value),
	}
}

func InternalTealValueToModel(value generated.TealValue) *model.TealValue {
	valueBytes, err := base64.RawStdEncoding.DecodeString(value.Bytes)
	if err != nil {
		panic(err)
	}
	return &model.TealValue{
		Bytes: valueBytes,
		Type:  value.Type,
		Uint:  value.Uint,
	}
}

func InternalApplicationLocalStatesToModel(localStates *[]generated.ApplicationLocalState) []model.ApplicationLocalState {
	if localStates == nil {
		return []model.ApplicationLocalState{}
	}

	converted := make([]model.ApplicationLocalState, len(*localStates))
	for i, localState := range *localStates {
		converted[i] = InternalApplicationLocalStateToModel(localState)
	}
	return converted
}

func InternalApplicationLocalStateToModel(localState generated.ApplicationLocalState) model.ApplicationLocalState {
	return model.ApplicationLocalState{
		ClosedOutAtRound: localState.ClosedOutAtRound,
		Deleted:          boolPtrOrDefault(localState.Deleted),
		ID:               localState.Id,
		KeyValue:         InternalTealKeyValueStoreToModel(localState.KeyValue),
		OptedInAtRound:   localState.OptedInAtRound,
		Schema:           InternalApplicationStateSchemaToModel(&localState.Schema),
	}
}

func InternalApplicationParamsToModel(params generated.ApplicationParams) *model.ApplicationParams {
	return &model.ApplicationParams{
		ApprovalProgram:   params.ApprovalProgram,
		ClearStateProgram: params.ClearStateProgram,
		Creator:           params.Creator,
		ExtraProgramPages: params.ExtraProgramPages,
		GlobalState:       InternalTealKeyValueStoreToModel(params.GlobalState),
		GlobalStateSchema: InternalApplicationStateSchemaToModel(params.GlobalStateSchema),
		LocalStateSchema:  InternalApplicationStateSchemaToModel(params.LocalStateSchema),
	}
}

func InternalApplicationsToModel(apps *[]generated.Application) []model.Application {
	if apps == nil {
		return []model.Application{}
	}

	converted := make([]model.Application, len(*apps))
	for i, app := range *apps {
		converted[i] = InternalApplicationToModel(app)
	}
	return converted
}

func InternalApplicationToModel(app generated.Application) model.Application {
	return model.Application{
		CreatedAtRound: app.CreatedAtRound,
		Deleted:        boolPtrOrDefault(app.Deleted),
		DeletedAtRound: app.DeletedAtRound,
		ID:             app.Id,
		Params:         InternalApplicationParamsToModel(app.Params),
	}
}

func InternalAssetHoldingsToModel(holdings *[]generated.AssetHolding) []model.AssetHolding {
	if holdings == nil {
		return []model.AssetHolding{}
	}

	converted := make([]model.AssetHolding, len(*holdings))
	for i, holding := range *holdings {
		converted[i] = InternalAssetHoldingToModel(holding)
	}
	return converted
}

func InternalAssetHoldingToModel(holding generated.AssetHolding) model.AssetHolding {
	return model.AssetHolding{
		Amount:          holding.Amount,
		ID:              holding.AssetId,
		Creator:         holding.Creator,
		Deleted:         boolPtrOrDefault(holding.Deleted),
		Frozen:          holding.IsFrozen,
		OptedInAtRound:  holding.OptedInAtRound,
		OptedOutAtRound: holding.OptedOutAtRound,
	}
}

func InternalAssetParamsToModel(params generated.AssetParams) *model.AssetParams {
	return &model.AssetParams{
		Clawback:      params.Clawback,
		Creator:       params.Creator,
		Decimals:      params.Decimals,
		DefaultFrozen: params.DefaultFrozen,
		Freeze:        params.Freeze,
		Manager:       params.Manager,
		MetadataHash:  byteslicePtrOrDefault(params.MetadataHash),
		Name:          params.Name,
		Reserve:       params.Reserve,
		Total:         params.Total,
		UnitName:      params.UnitName,
		URL:           params.Url,
	}
}

func InternalAssetsToModel(assets *[]generated.Asset) []model.Asset {
	if assets == nil {
		return []model.Asset{}
	}

	converted := make([]model.Asset, len(*assets))
	for i, asset := range *assets {
		converted[i] = InternalAssetToModel(asset)
	}
	return converted
}

func InternalAssetToModel(asset generated.Asset) model.Asset {
	return model.Asset{
		CreatedAtRound:   asset.CreatedAtRound,
		Deleted:          boolPtrOrDefault(asset.Deleted),
		DestroyedAtRound: asset.DestroyedAtRound,
		ID:               asset.Index,
		Params:           InternalAssetParamsToModel(asset.Params),
	}
}
