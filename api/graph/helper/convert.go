package helper

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/algorand/go-algorand-sdk/crypto"
	sdk_types "github.com/algorand/go-algorand-sdk/types"

	"github.com/algorand/indexer/api/generated/v2"
	"github.com/algorand/indexer/api/graph/model"
	"github.com/algorand/indexer/types"
)

func boolPtr(b bool) *bool {
	return &b
}

func boolOrDefault(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func uint64Ptr(i uint64) *uint64 {
	return &i
}

func uint64OrDefault(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}

func uint64SliceOrDefault(p *[]uint64) []uint64 {
	if p == nil {
		return []uint64{}
	}
	return *p
}

func bytePtr(x []byte) *[]byte {
	if len(x) == 0 {
		return nil
	}

	// Don't return if it's all zero.
	for _, v := range x {
		if v != 0 {
			return &x
		}
	}

	return nil
}

func byteSliceOrDefault(p *[]byte) []byte {
	if p == nil {
		return []byte{}
	}
	return *p
}

func strPtr(s string) *string {
	return &s
}

func stringSliceOrDefault(p *[]string) []string {
	if p == nil {
		return []string{}
	}
	return *p
}

func addrPtr(x sdk_types.Address) *string {
	if x.IsZero() {
		return nil
	}
	out := new(string)
	*out = x.String()
	return out
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

func ModalSigTypeToInternal(sigType *model.SigType) *string {
	if sigType == nil {
		return nil
	}

	result := strings.ToLower(sigType.String())
	return &result
}

func InternalAccountStatusToModel(status *string) *model.AccountStatus {
	if status == nil {
		return nil
	}

	var result model.AccountStatus

	switch *status {
	case "Offline":
		result = model.AccountStatusOffline
	case "Online":
		result = model.AccountStatusOnline
	case "NotParticipating":
		result = model.AccountStatusNotParticipating
	default:
		panic(fmt.Errorf("Unexpected account status: %s", *status))
	}

	return &result
}

func InternalTxTypeToModel(txType *string) *model.TxType {
	if txType == nil {
		return nil
	}

	var result model.TxType

	switch *txType {
	case "pay":
		result = model.TxTypePay
	case "keyreg":
		result = model.TxTypeKeyreg
	case "acfg":
		result = model.TxTypeAcfg
	case "axfer":
		result = model.TxTypeAxfer
	case "afrz":
		result = model.TxTypeAfrz
	case "appl":
		result = model.TxTypeAppl
	default:
		panic(fmt.Errorf("Unexpected transaction type: %s", *txType))
	}

	return &result
}

func ModalTxTypeToInternal(txType *model.TxType) *string {
	if txType == nil {
		return nil
	}

	result := strings.ToLower(txType.String())
	return &result
}

func InternalSDKOnCompletionToModel(oc sdk_types.OnCompletion) model.OnCompletion {
	switch oc {
	case sdk_types.NoOpOC:
		return model.OnCompletionNoop
	case sdk_types.OptInOC:
		return model.OnCompletionOptin
	case sdk_types.CloseOutOC:
		return model.OnCompletionCloseout
	case sdk_types.ClearStateOC:
		return model.OnCompletionClear
	case sdk_types.UpdateApplicationOC:
		return model.OnCompletionUpdate
	case sdk_types.DeleteApplicationOC:
		return model.OnCompletionDelete
	default:
		panic(fmt.Errorf("Unexpected OnCompletion value: %d", oc))
	}
}

func InternalOnCompletionToModel(oc generated.OnCompletion) model.OnCompletion {
	switch oc {
	case "noop":
		return model.OnCompletionNoop
	case "optin":
		return model.OnCompletionOptin
	case "closeout":
		return model.OnCompletionCloseout
	case "clear":
		return model.OnCompletionClear
	case "update":
		return model.OnCompletionUpdate
	case "delete":
		return model.OnCompletionDelete
	default:
		panic(fmt.Errorf("Unexpected OnCompletion value: %s", oc))
	}
}

func ModalAddressRoleToInternal(role *model.AddressRole) *string {
	if role == nil {
		return nil
	}

	var result string

	switch *role {
	case model.AddressRoleSender:
		result = "sender"
	case model.AddressRoleReceiver:
		result = "receiver"
	case model.AddressRoleFreezeTarget:
		result = "freeze-target"
	default:
		panic(fmt.Errorf("Unexpected AddressRole value: %s", *role))
	}

	return &result
}

func InternalDeltaActionToModel(action uint64) model.DeltaAction {
	switch action {
	case 1:
		return model.DeltaActionSetBytes
	case 2:
		return model.DeltaActionSetUINt
	case 3:
		return model.DeltaActionDelete
	default:
		panic(fmt.Errorf("Unexpected action: %d", action))
	}
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
		AppsTotalExtraPages:         uint64OrDefault(account.AppsTotalExtraPages),
		AppsTotalSchema:             InternalApplicationStateSchemaToModel(account.AppsTotalSchema),
		Assets:                      InternalAssetHoldingsToModel(account.Assets),
		AuthAddr:                    account.AuthAddr,
		ClosedAtRound:               account.ClosedAtRound,
		CreatedApps:                 InternalApplicationsToModel(account.CreatedApps),
		CreatedAssets:               InternalAssetsToModel(account.CreatedAssets),
		CreatedAtRound:              account.CreatedAtRound,
		Deleted:                     boolOrDefault(account.Deleted),
		Participation:               InternalAccountParticipationToModel(account.Participation),
		PendingRewards:              account.PendingRewards,
		RewardBase:                  account.RewardBase,
		Rewards:                     account.Rewards,
		Round:                       account.Round,
		SigType:                     InternalSigTypeToModel(account.SigType),
		Status:                      *InternalAccountStatusToModel(&account.Status),
	}
}

func InternalApplicationStateSchemaToModel(schema *generated.ApplicationStateSchema) *model.ApplicationStateSchema {
	if schema == nil {
		return &model.ApplicationStateSchema{}
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
		Deleted:          boolOrDefault(localState.Deleted),
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
		Creator:           *params.Creator,
		ExtraProgramPages: uint64OrDefault(params.ExtraProgramPages),
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
		converted[i] = *InternalApplicationToModel(app)
	}
	return converted
}

func InternalApplicationToModel(app generated.Application) *model.Application {
	return &model.Application{
		CreatedAtRound: app.CreatedAtRound,
		Deleted:        boolOrDefault(app.Deleted),
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
		Deleted:         boolOrDefault(holding.Deleted),
		Frozen:          holding.IsFrozen,
		OptedInAtRound:  holding.OptedInAtRound,
		OptedOutAtRound: holding.OptedOutAtRound,
	}
}

func InternalMiniAssetHoldingsToModel(holdings *[]generated.MiniAssetHolding) []model.MiniAssetHolding {
	if holdings == nil {
		return []model.MiniAssetHolding{}
	}

	converted := make([]model.MiniAssetHolding, len(*holdings))
	for i, holding := range *holdings {
		converted[i] = InternalMiniAssetHoldingToModel(holding)
	}
	return converted
}

func InternalMiniAssetHoldingToModel(holding generated.MiniAssetHolding) model.MiniAssetHolding {
	return model.MiniAssetHolding{
		Address:         holding.Address,
		Amount:          holding.Amount,
		Deleted:         boolOrDefault(holding.Deleted),
		Frozen:          holding.IsFrozen,
		OptedInAtRound:  holding.OptedInAtRound,
		OptedOutAtRound: holding.OptedOutAtRound,
	}
}

func InternalAssetParamsToModel(params *generated.AssetParams) *model.AssetParams {
	if params == nil {
		return nil
	}

	return &model.AssetParams{
		Clawback:      params.Clawback,
		Creator:       params.Creator,
		Decimals:      params.Decimals,
		DefaultFrozen: boolOrDefault(params.DefaultFrozen),
		Freeze:        params.Freeze,
		Manager:       params.Manager,
		MetadataHash:  byteSliceOrDefault(params.MetadataHash),
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
		converted[i] = *InternalAssetToModel(asset)
	}
	return converted
}

func InternalAssetToModel(asset generated.Asset) *model.Asset {
	return &model.Asset{
		CreatedAtRound:   asset.CreatedAtRound,
		Deleted:          boolOrDefault(asset.Deleted),
		DestroyedAtRound: asset.DestroyedAtRound,
		ID:               asset.Index,
		Params:           InternalAssetParamsToModel(&asset.Params),
	}
}

func InternalBlockRewardsToModel(rewards *generated.BlockRewards) *model.BlockRewards {
	if rewards == nil {
		return nil
	}

	return &model.BlockRewards{
		FeeSink:                 rewards.FeeSink,
		RewardsCalculationRound: rewards.RewardsCalculationRound,
		RewardsLevel:            rewards.RewardsLevel,
		RewardsPool:             rewards.RewardsPool,
		RewardsRate:             rewards.RewardsRate,
		RewardsResidue:          rewards.RewardsResidue,
	}
}

func InternalBlockUpgradeStateToModel(state *generated.BlockUpgradeState) *model.BlockUpgradeState {
	if state == nil {
		return nil
	}

	return &model.BlockUpgradeState{
		CurrentProtocol:        state.CurrentProtocol,
		NextProtocol:           state.NextProtocol,
		NextProtocolApprovals:  state.NextProtocolApprovals,
		NextProtocolSwitchOn:   state.NextProtocolSwitchOn,
		NextProtocolVoteBefore: state.NextProtocolVoteBefore,
	}
}

func InternalBlockUpgradeVoteToModel(vote *generated.BlockUpgradeVote) *model.BlockUpgradeVote {
	if vote == nil {
		return nil
	}

	return &model.BlockUpgradeVote{
		UpgradeApprove: vote.UpgradeApprove,
		UpgradeDelay:   vote.UpgradeDelay,
		UpgradePropose: vote.UpgradePropose,
	}
}

func InternalBlockHeaderAndTxnsToModel(header types.BlockHeader, txns []types.SignedTxnWithAD, createdPrimitives map[int]uint64) *model.Block {
	rewards := model.BlockRewards{
		FeeSink:                 header.FeeSink.String(),
		RewardsCalculationRound: uint64(header.RewardsRecalculationRound),
		RewardsLevel:            header.RewardsLevel,
		RewardsPool:             header.RewardsPool.String(),
		RewardsRate:             header.RewardsRate,
		RewardsResidue:          header.RewardsResidue,
	}

	upgradeState := model.BlockUpgradeState{
		CurrentProtocol:        string(header.CurrentProtocol),
		NextProtocol:           strPtr(string(header.NextProtocol)),
		NextProtocolApprovals:  uint64Ptr(header.NextProtocolApprovals),
		NextProtocolSwitchOn:   uint64Ptr(uint64(header.NextProtocolSwitchOn)),
		NextProtocolVoteBefore: uint64Ptr(uint64(header.NextProtocolVoteBefore)),
	}

	upgradeVote := model.BlockUpgradeVote{
		UpgradeApprove: boolPtr(header.UpgradeApprove),
		UpgradeDelay:   uint64Ptr(uint64(header.UpgradeDelay)),
		UpgradePropose: strPtr(string(header.UpgradePropose)),
	}

	ret := model.Block{
		GenesisHash:       header.GenesisHash[:],
		GenesisID:         header.GenesisID,
		PreviousBlockHash: header.Branch[:],
		Rewards:           &rewards,
		Round:             uint64(header.Round),
		Seed:              header.Seed[:],
		Timestamp:         uint64(header.TimeStamp),
		TransactionsRoot:  header.TxnRoot[:],
		TxnCounter:        uint64Ptr(header.TxnCounter),
		UpgradeState:      &upgradeState,
		UpgradeVote:       &upgradeVote,
	}

	results := make([]model.Transaction, len(txns))
	for i, txn := range txns {
		created := createdPrimitives[i]
		results[i] = InternalSignedTxnWithADToModel(txn, created, uint64(header.Round), i, uint64(header.TimeStamp))
	}

	ret.Transactions = results
	return &ret
}

func msigToTransactionMsig(msig sdk_types.MultisigSig) *model.TransactionSignatureMultisig {
	if msig.Blank() {
		return nil
	}

	subsigs := make([]model.TransactionSignatureMultisigSubsignature, 0)
	for _, subsig := range msig.Subsigs {
		subsigs = append(subsigs, model.TransactionSignatureMultisigSubsignature{
			PublicKey: subsig.Key[:],
			Signature: subsig.Sig[:],
		})
	}

	ret := model.TransactionSignatureMultisig{
		Subsignature: subsigs,
		Threshold:    uint64Ptr(uint64(msig.Threshold)),
		Version:      uint64Ptr(uint64(msig.Version)),
	}
	return &ret
}

func lsigToTransactionLsig(lsig sdk_types.LogicSig) *model.TransactionSignatureLogicsig {
	if lsig.Blank() {
		return nil
	}

	ret := model.TransactionSignatureLogicsig{
		Args:              lsig.Args,
		Logic:             lsig.Logic,
		MultisigSignature: msigToTransactionMsig(lsig.Msig),
		Signature:         lsig.Sig[:],
	}

	return &ret
}

func stateDeltaToStateDelta(d types.StateDelta) []model.EvalDeltaKeyValue {
	if len(d) == 0 {
		return nil
	}
	var delta []model.EvalDeltaKeyValue
	keys := make([]string, 0)
	for k := range d {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := d[k]
		delta = append(delta, model.EvalDeltaKeyValue{
			Key: []byte(k),
			Value: &model.EvalDelta{
				Action: InternalDeltaActionToModel(uint64(v.Action)),
				Bytes:  v.Bytes,
				Uint:   uint64Ptr(v.Uint),
			},
		})
	}
	return delta
}

func InternalSignedTxnWithADToModel(stxn types.SignedTxnWithAD, createdPrimitive uint64, round uint64, intra int, roundTime uint64) model.Transaction {
	var payment *model.TransactionPayment
	var keyreg *model.TransactionKeyreg
	var assetConfig *model.TransactionAssetConfig
	var assetFreeze *model.TransactionAssetFreeze
	var assetTransfer *model.TransactionAssetTransfer
	var application *model.TransactionApplication

	switch stxn.Txn.Type {
	case sdk_types.PaymentTx:
		p := model.TransactionPayment{
			CloseAmount:      uint64Ptr(uint64(stxn.ApplyData.ClosingAmount)),
			CloseRemainderTo: addrPtr(stxn.Txn.CloseRemainderTo),
			Receiver:         stxn.Txn.Receiver.String(),
			Amount:           uint64(stxn.Txn.Amount),
		}
		payment = &p
	case sdk_types.KeyRegistrationTx:
		k := model.TransactionKeyreg{
			NonParticipation:          stxn.Txn.Nonparticipation,
			SelectionParticipationKey: stxn.Txn.SelectionPK[:],
			VoteFirstValid:            uint64Ptr(uint64(stxn.Txn.VoteFirst)),
			VoteLastValid:             uint64Ptr(uint64(stxn.Txn.VoteLast)),
			VoteKeyDilution:           uint64Ptr(stxn.Txn.VoteKeyDilution),
			VoteParticipationKey:      stxn.Txn.VotePK[:],
		}
		keyreg = &k
	case sdk_types.AssetConfigTx:
		assetParams := model.AssetParams{
			Clawback:      addrPtr(stxn.Txn.AssetParams.Clawback),
			Creator:       stxn.Txn.Sender.String(),
			Decimals:      uint64(stxn.Txn.AssetParams.Decimals),
			DefaultFrozen: stxn.Txn.AssetParams.DefaultFrozen,
			Freeze:        addrPtr(stxn.Txn.AssetParams.Freeze),
			Manager:       addrPtr(stxn.Txn.AssetParams.Manager),
			MetadataHash:  stxn.Txn.AssetParams.MetadataHash[:],
			Name:          strPtr(stxn.Txn.AssetParams.AssetName),
			Reserve:       addrPtr(stxn.Txn.AssetParams.Reserve),
			Total:         stxn.Txn.AssetParams.Total,
			UnitName:      strPtr(stxn.Txn.AssetParams.UnitName),
			URL:           strPtr(stxn.Txn.AssetParams.URL),
		}
		config := model.TransactionAssetConfig{
			AssetID: uint64Ptr(uint64(stxn.Txn.ConfigAsset)),
			Params:  &assetParams,
		}
		assetConfig = &config
	case sdk_types.AssetTransferTx:
		t := model.TransactionAssetTransfer{
			Amount:      stxn.Txn.AssetAmount,
			AssetID:     uint64(stxn.Txn.XferAsset),
			CloseTo:     addrPtr(stxn.Txn.AssetCloseTo),
			Receiver:    stxn.Txn.AssetReceiver.String(),
			Sender:      addrPtr(stxn.Txn.AssetSender),
			CloseAmount: uint64Ptr(stxn.ApplyData.AssetClosingAmount),
		}
		assetTransfer = &t
	case sdk_types.AssetFreezeTx:
		f := model.TransactionAssetFreeze{
			Address:         stxn.Txn.FreezeAccount.String(),
			AssetID:         uint64(stxn.Txn.FreezeAsset),
			NewFreezeStatus: stxn.Txn.AssetFrozen,
		}
		assetFreeze = &f
	case sdk_types.ApplicationCallTx:
		accts := make([]string, len(stxn.Txn.Accounts))
		for _, v := range stxn.Txn.Accounts {
			accts = append(accts, v.String())
		}

		apps := make([]uint64, len(stxn.Txn.ForeignApps))
		for _, v := range stxn.Txn.ForeignApps {
			apps = append(apps, uint64(v))
		}

		assets := make([]uint64, len(stxn.Txn.ForeignAssets))
		for _, v := range stxn.Txn.ForeignAssets {
			assets = append(assets, uint64(v))
		}

		a := model.TransactionApplication{
			Accounts:          accts,
			ApplicationArgs:   stxn.Txn.ApplicationArgs,
			ApplicationID:     uint64(stxn.Txn.ApplicationID),
			ApprovalProgram:   stxn.Txn.ApprovalProgram,
			ClearStateProgram: stxn.Txn.ClearStateProgram,
			ForeignApps:       apps,
			ForeignAssets:     assets,
			GlobalStateSchema: &model.StateSchema{
				NumByteSlice: stxn.Txn.GlobalStateSchema.NumByteSlice,
				NumUint:      stxn.Txn.GlobalStateSchema.NumUint,
			},
			LocalStateSchema: &model.StateSchema{
				NumByteSlice: stxn.Txn.LocalStateSchema.NumByteSlice,
				NumUint:      stxn.Txn.LocalStateSchema.NumUint,
			},
			OnCompletion:      InternalSDKOnCompletionToModel(stxn.Txn.OnCompletion),
			ExtraProgramPages: uint64Ptr(uint64(stxn.Txn.ExtraProgramPages)),
		}

		application = &a
	}

	sig := model.TransactionSignature{
		Logicsig: lsigToTransactionLsig(stxn.Lsig),
		Multisig: msigToTransactionMsig(stxn.Msig),
		Sig:      stxn.Sig[:],
	}

	var localStateDelta []model.AccountStateDelta
	type tuple struct {
		key     uint64
		address types.Address
	}
	if len(stxn.ApplyData.EvalDelta.LocalDeltas) > 0 {
		keys := make([]tuple, 0)
		for k := range stxn.ApplyData.EvalDelta.LocalDeltas {
			if k == 0 {
				keys = append(keys, tuple{
					key:     0,
					address: stxn.Txn.Sender,
				})
			} else {
				addr := types.Address{}
				copy(addr[:], stxn.Txn.Accounts[k-1][:])
				keys = append(keys, tuple{
					key:     k,
					address: addr,
				})
			}
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i].key < keys[j].key })
		localStateDelta := make([]model.AccountStateDelta, 0)
		for _, k := range keys {
			v := stxn.ApplyData.EvalDelta.LocalDeltas[k.key]
			delta := stateDeltaToStateDelta(v)
			if delta != nil {
				localStateDelta = append(localStateDelta, model.AccountStateDelta{
					Address: k.address.String(),
					Delta:   delta,
				})
			}
		}
	}

	txn := model.Transaction{
		ApplicationTransaction:   application,
		AssetConfigTransaction:   assetConfig,
		AssetFreezeTransaction:   assetFreeze,
		AssetTransferTransaction: assetTransfer,
		PaymentTransaction:       payment,
		KeyregTransaction:        keyreg,
		ClosingAmount:            uint64Ptr(uint64(stxn.ClosingAmount)),
		ConfirmedRound:           uint64Ptr(round),
		IntraRoundOffset:         uint64Ptr(uint64(intra)),
		RoundTime:                uint64Ptr(roundTime),
		Fee:                      uint64(stxn.Txn.Fee),
		FirstValid:               uint64(stxn.Txn.FirstValid),
		GenesisHash:              stxn.SignedTxn.Txn.GenesisHash[:],
		GenesisID:                strPtr(stxn.SignedTxn.Txn.GenesisID),
		Group:                    stxn.Txn.Group[:],
		LastValid:                uint64(stxn.Txn.LastValid),
		Lease:                    stxn.Txn.Lease[:],
		Note:                     stxn.Txn.Note[:],
		Sender:                   stxn.Txn.Sender.String(),
		ReceiverRewards:          uint64Ptr(uint64(stxn.ReceiverRewards)),
		CloseRewards:             uint64Ptr(uint64(stxn.CloseRewards)),
		SenderRewards:            uint64Ptr(uint64(stxn.SenderRewards)),
		TxType:                   *InternalTxTypeToModel(strPtr(string(stxn.Txn.Type))),
		Signature:                &sig,
		ID:                       crypto.TransactionIDString(stxn.Txn),
		RekeyTo:                  addrPtr(stxn.Txn.RekeyTo),
		GlobalStateDelta:         stateDeltaToStateDelta(stxn.EvalDelta.GlobalDelta),
		LocalStateDelta:          localStateDelta,
	}

	if stxn.Txn.Type == sdk_types.AssetConfigTx {
		if txn.AssetConfigTransaction != nil && txn.AssetConfigTransaction.AssetID != nil && *txn.AssetConfigTransaction.AssetID == 0 {
			txn.CreatedAssetID = uint64Ptr(createdPrimitive)
		}
	}

	if stxn.Txn.Type == sdk_types.ApplicationCallTx {
		if txn.ApplicationTransaction != nil && txn.ApplicationTransaction.ApplicationID == 0 {
			txn.CreatedApplicationID = uint64Ptr(createdPrimitive)
		}
	}

	return txn
}

func InternalBlockToModel(blk generated.Block) *model.Block {
	return &model.Block{
		GenesisHash:       blk.GenesisHash,
		GenesisID:         blk.GenesisId,
		PreviousBlockHash: blk.PreviousBlockHash,
		Rewards:           InternalBlockRewardsToModel(blk.Rewards),
		Round:             blk.Round,
		Seed:              blk.Seed,
		Timestamp:         blk.Timestamp,
		Transactions:      InternalTransactionsToModel(blk.Transactions),
		TransactionsRoot:  blk.TransactionsRoot,
		TxnCounter:        blk.TxnCounter,
		UpgradeState:      InternalBlockUpgradeStateToModel(blk.UpgradeState),
		UpgradeVote:       InternalBlockUpgradeVoteToModel(blk.UpgradeVote),
	}
}

func InternalTransactionsToModel(txns *[]generated.Transaction) []model.Transaction {
	if txns == nil {
		return []model.Transaction{}
	}

	converted := make([]model.Transaction, len(*txns))
	for i, txn := range *txns {
		converted[i] = *InternalTransactionToModel(txn)
	}
	return converted
}

func InternalTransactionToModel(txn generated.Transaction) *model.Transaction {
	return &model.Transaction{
		ApplicationTransaction:   InternalTransactionApplicationToModel(txn.ApplicationTransaction),
		AssetConfigTransaction:   InternalTransactionAssetConfigToModel(txn.AssetConfigTransaction),
		AssetFreezeTransaction:   InternalTransactionAssetFreezeToModel(txn.AssetFreezeTransaction),
		AssetTransferTransaction: InternalTransactionAssetTransferToModel(txn.AssetTransferTransaction),
		AuthAddr:                 txn.AuthAddr,
		CloseRewards:             txn.CloseRewards,
		ClosingAmount:            txn.ClosingAmount,
		ConfirmedRound:           txn.ConfirmedRound,
		CreatedApplicationID:     txn.CreatedApplicationIndex,
		CreatedAssetID:           txn.CreatedAssetIndex,
		Fee:                      txn.Fee,
		FirstValid:               txn.FirstValid,
		GenesisHash:              byteSliceOrDefault(txn.GenesisHash),
		GenesisID:                txn.GenesisId,
		GlobalStateDelta:         InternalEvalDeltaKeyValuesToModel(txn.GlobalStateDelta),
		Group:                    byteSliceOrDefault(txn.Group),
		ID:                       txn.Id,
		IntraRoundOffset:         txn.IntraRoundOffset,
		KeyregTransaction:        InternalTransactionKeyregToModel(txn.KeyregTransaction),
		LastValid:                txn.LastValid,
		Lease:                    byteSliceOrDefault(txn.Lease),
		LocalStateDelta:          InternalAccountStateDeltasToModel(txn.LocalStateDelta),
		Note:                     byteSliceOrDefault(txn.Note),
		PaymentTransaction:       InternalTransactionPaymentToModel(txn.PaymentTransaction),
		ReceiverRewards:          txn.ReceiverRewards,
		RekeyTo:                  txn.RekeyTo,
		RoundTime:                txn.RoundTime,
		Sender:                   txn.Sender,
		SenderRewards:            txn.SenderRewards,
		Signature:                InternalTransactionSignatureToModel(&txn.Signature),
		TxType:                   *InternalTxTypeToModel(&txn.TxType),
	}
}

func InternalTransactionSignatureToModel(sig *generated.TransactionSignature) *model.TransactionSignature {
	if sig == nil {
		return nil
	}

	return &model.TransactionSignature{
		Logicsig: nil,
		Multisig: nil,
		Sig:      byteSliceOrDefault(sig.Sig),
	}
}

func InternalTransactionSignatureLogicsigToModel(lsig *generated.TransactionSignatureLogicsig) *model.TransactionSignatureLogicsig {
	if lsig == nil {
		return nil
	}

	var args [][]byte
	if lsig.Args != nil {
		args = make([][]byte, len(*lsig.Args))
		for i, argStr := range *lsig.Args {
			b64, err := base64.RawStdEncoding.DecodeString(argStr)
			if err != nil {
				panic(err)
			}
			args[i] = b64
		}
	}

	return &model.TransactionSignatureLogicsig{
		Args:              args,
		Logic:             lsig.Logic,
		MultisigSignature: InternalTransactionSignatureMultisigToModel(lsig.MultisigSignature),
		Signature:         byteSliceOrDefault(lsig.Signature),
	}
}

func InternalTransactionSignatureMultisigToModel(msig *generated.TransactionSignatureMultisig) *model.TransactionSignatureMultisig {
	if msig == nil {
		return nil
	}

	return &model.TransactionSignatureMultisig{
		Subsignature: InternalTransactionSignatureMultisigSubsignaturesToModel(msig.Subsignature),
		Threshold:    msig.Threshold,
		Version:      msig.Version,
	}
}

func InternalTransactionSignatureMultisigSubsignaturesToModel(subsigs *[]generated.TransactionSignatureMultisigSubsignature) []model.TransactionSignatureMultisigSubsignature {
	if subsigs == nil {
		return []model.TransactionSignatureMultisigSubsignature{}
	}

	converted := make([]model.TransactionSignatureMultisigSubsignature, len(*subsigs))
	for i, subsig := range *subsigs {
		converted[i] = InternalTransactionSignatureMultisigSubsignatureToModel(subsig)
	}
	return converted
}

func InternalTransactionSignatureMultisigSubsignatureToModel(subsig generated.TransactionSignatureMultisigSubsignature) model.TransactionSignatureMultisigSubsignature {
	return model.TransactionSignatureMultisigSubsignature{
		PublicKey: byteSliceOrDefault(subsig.PublicKey),
		Signature: byteSliceOrDefault(subsig.Signature),
	}
}

func InternalEvalDeltaKeyValuesToModel(deltas *generated.StateDelta) []model.EvalDeltaKeyValue {
	if deltas == nil {
		return []model.EvalDeltaKeyValue{}
	}

	converted := make([]model.EvalDeltaKeyValue, len(*deltas))
	for i, delta := range *deltas {
		converted[i] = InternalEvalDeltaKeyValueToModel(delta)
	}
	return converted
}

func InternalEvalDeltaKeyValueToModel(delta generated.EvalDeltaKeyValue) model.EvalDeltaKeyValue {
	key, err := base64.RawStdEncoding.DecodeString(delta.Key)
	if err != nil {
		panic(err)
	}

	return model.EvalDeltaKeyValue{
		Key:   key,
		Value: InternalEvalDeltaToModel(delta.Value),
	}
}

func InternalEvalDeltaToModel(delta generated.EvalDelta) *model.EvalDelta {
	var bytes []byte
	if delta.Bytes != nil {
		var err error
		bytes, err = base64.RawStdEncoding.DecodeString(*delta.Bytes)
		if err != nil {
			panic(err)
		}
	}

	return &model.EvalDelta{
		Action: InternalDeltaActionToModel(delta.Action),
		Bytes:  bytes,
		Uint:   delta.Uint,
	}
}

func InternalAccountStateDeltasToModel(deltas *[]generated.AccountStateDelta) []model.AccountStateDelta {
	if deltas == nil {
		return []model.AccountStateDelta{}
	}

	converted := make([]model.AccountStateDelta, len(*deltas))
	for i, delta := range *deltas {
		converted[i] = InternalAccountStateDeltaToModel(delta)
	}
	return converted
}

func InternalAccountStateDeltaToModel(delta generated.AccountStateDelta) model.AccountStateDelta {
	return model.AccountStateDelta{
		Address: delta.Address,
		Delta:   InternalEvalDeltaKeyValuesToModel(&delta.Delta),
	}
}

func InternalTransactionPaymentToModel(txn *generated.TransactionPayment) *model.TransactionPayment {
	if txn == nil {
		return nil
	}

	return &model.TransactionPayment{
		Amount:           txn.Amount,
		CloseAmount:      txn.CloseAmount,
		CloseRemainderTo: txn.CloseRemainderTo,
		Receiver:         txn.Receiver,
	}
}

func InternalTransactionKeyregToModel(txn *generated.TransactionKeyreg) *model.TransactionKeyreg {
	if txn == nil {
		return nil
	}

	return &model.TransactionKeyreg{
		NonParticipation:          boolOrDefault(txn.NonParticipation),
		SelectionParticipationKey: byteSliceOrDefault(txn.SelectionParticipationKey),
		VoteFirstValid:            txn.VoteFirstValid,
		VoteKeyDilution:           txn.VoteKeyDilution,
		VoteLastValid:             txn.VoteLastValid,
		VoteParticipationKey:      byteSliceOrDefault(txn.VoteParticipationKey),
	}
}

func InternalTransactionAssetTransferToModel(txn *generated.TransactionAssetTransfer) *model.TransactionAssetTransfer {
	if txn == nil {
		return nil
	}

	return &model.TransactionAssetTransfer{
		Amount:      txn.Amount,
		AssetID:     txn.AssetId,
		CloseAmount: txn.CloseAmount,
		CloseTo:     txn.CloseTo,
		Receiver:    txn.Receiver,
		Sender:      txn.Sender,
	}
}

func InternalTransactionAssetFreezeToModel(txn *generated.TransactionAssetFreeze) *model.TransactionAssetFreeze {
	if txn == nil {
		return nil
	}

	return &model.TransactionAssetFreeze{
		Address:         txn.Address,
		AssetID:         txn.AssetId,
		NewFreezeStatus: txn.NewFreezeStatus,
	}
}

func InternalTransactionAssetConfigToModel(txn *generated.TransactionAssetConfig) *model.TransactionAssetConfig {
	if txn == nil {
		return nil
	}

	return &model.TransactionAssetConfig{
		AssetID: txn.AssetId,
		Params:  InternalAssetParamsToModel(txn.Params),
	}
}

func InternalTransactionApplicationToModel(txn *generated.TransactionApplication) *model.TransactionApplication {
	if txn == nil {
		return nil
	}

	var args [][]byte
	if txn.ApplicationArgs != nil {
		args = make([][]byte, len(*txn.ApplicationArgs))
		for i, argStr := range *txn.ApplicationArgs {
			b64, err := base64.RawStdEncoding.DecodeString(argStr)
			if err != nil {
				panic(err)
			}
			args[i] = b64
		}
	}

	return &model.TransactionApplication{
		Accounts:          stringSliceOrDefault(txn.Accounts),
		ApplicationArgs:   args,
		ApplicationID:     txn.ApplicationId,
		ApprovalProgram:   byteSliceOrDefault(txn.ApprovalProgram),
		ClearStateProgram: byteSliceOrDefault(txn.ClearStateProgram),
		ExtraProgramPages: txn.ExtraProgramPages,
		ForeignApps:       uint64SliceOrDefault(txn.ForeignApps),
		ForeignAssets:     uint64SliceOrDefault(txn.ForeignAssets),
		GlobalStateSchema: InternalStateSchemaToModel(txn.GlobalStateSchema),
		LocalStateSchema:  InternalStateSchemaToModel(txn.LocalStateSchema),
		OnCompletion:      InternalOnCompletionToModel(txn.OnCompletion),
	}
}

func InternalStateSchemaToModel(schema *generated.StateSchema) *model.StateSchema {
	if schema == nil {
		return nil
	}

	return &model.StateSchema{
		NumByteSlice: schema.NumByteSlice,
		NumUint:      schema.NumUint,
	}
}
