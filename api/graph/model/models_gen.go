// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

// Account information at a given round.
//
// Definition:
// data/basics/userBalance.go : AccountData
type Account struct {
	// the account public key
	Address string `json:"address"`
	// \[algo\] total number of MicroAlgos in the account
	Amount uint64 `json:"amount"`
	// specifies the amount of MicroAlgos in the account, without the pending rewards.
	AmountWithoutPendingRewards uint64 `json:"amountWithoutPendingRewards"`
	// \[appl\] applications local data stored in this account.
	//
	// Note the raw object uses `map[int] -> AppLocalState` for this type.
	AppsLocalState []ApplicationLocalState `json:"appsLocalState"`
	// \[teap\] the sum of all extra application program pages for this account.
	AppsTotalExtraPages uint64 `json:"appsTotalExtraPages"`
	// Specifies maximums on the number of each type that may be stored.
	AppsTotalSchema *ApplicationStateSchema `json:"appsTotalSchema"`
	// \[asset\] assets held by this account.
	//
	// Note the raw object uses `map[int] -> AssetHolding` for this type.
	Assets []AssetHolding `json:"assets"`
	// \[spend\] the address against which signing should be checked. If empty, the address of the current account is used. This field can be updated in any transaction by setting the RekeyTo field.
	AuthAddr *string `json:"authAddr"`
	// Round during which this account was most recently closed.
	ClosedAtRound *uint64 `json:"closedAtRound"`
	// \[appp\] parameters of applications created by this account including app global data.
	//
	// Note: the raw account uses `map[int] -> AppParams` for this type.
	CreatedApps []Application `json:"createdApps"`
	// \[apar\] parameters of assets created by this account.
	//
	// Note: the raw account uses `map[int] -> Asset` for this type.
	CreatedAssets []Asset `json:"createdAssets"`
	// Round during which this account first appeared in a transaction.
	CreatedAtRound *uint64 `json:"createdAtRound"`
	// Whether or not this account is currently closed.
	Deleted bool `json:"deleted"`
	// AccountParticipation describes the parameters used by this account in consensus protocol.
	Participation *AccountParticipation `json:"participation"`
	// amount of MicroAlgos of pending rewards in this account.
	PendingRewards uint64 `json:"pendingRewards"`
	// \[ebase\] used as part of the rewards computation. Only applicable to accounts which are participating.
	RewardBase *uint64 `json:"rewardBase"`
	// \[ern\] total rewards of MicroAlgos the account has received, including pending rewards.
	Rewards uint64 `json:"rewards"`
	// The round for which this information is relevant.
	Round uint64 `json:"round"`
	// Indicates what type of signature is used by this account, must be one of:
	// * sig
	// * msig
	// * lsig
	SigType *SigType `json:"sigType"`
	// \[onl\] delegation status of the account's MicroAlgos
	// * Offline - indicates that the associated account is delegated.
	// *  Online  - indicates that the associated account used as part of the delegation pool.
	// *   NotParticipating - indicates that the associated account is neither a delegator nor a delegate.
	Status AccountStatus `json:"status"`
}

// AccountParticipation describes the parameters used by this account in consensus protocol.
type AccountParticipation struct {
	// \[sel\] Selection public key (if any) currently registered for this round.
	SelectionParticipationKey []byte `json:"selectionParticipationKey"`
	// \[voteFst\] First round for which this participation is valid.
	VoteFirstValid uint64 `json:"voteFirstValid"`
	// \[voteKD\] Number of subkeys in each batch of participation keys.
	VoteKeyDilution uint64 `json:"voteKeyDilution"`
	// \[voteLst\] Last round for which this participation is valid.
	VoteLastValid uint64 `json:"voteLastValid"`
	// \[vote\] root participation public key (if any) currently registered for this round.
	VoteParticipationKey []byte `json:"voteParticipationKey"`
}

type AccountResponse struct {
	// Account information at a given round.
	//
	// Definition:
	// data/basics/userBalance.go : AccountData
	Account *Account `json:"account"`
	// Round at which the results were computed.
	CurrentRound uint64 `json:"currentRound"`
}

// Application state delta.
type AccountStateDelta struct {
	Address string `json:"address"`
	// Application state delta.
	Delta []EvalDeltaKeyValue `json:"delta"`
}

type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
	// Round at which the results were computed.
	CurrentRound uint64 `json:"currentRound"`
	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken *string `json:"nextToken"`
}

// Application index and its parameters
type Application struct {
	// Round when this application was created.
	CreatedAtRound *uint64 `json:"createdAtRound"`
	// Whether or not this application is currently deleted.
	Deleted bool `json:"deleted"`
	// Round when this application was deleted.
	DeletedAtRound *uint64 `json:"deletedAtRound"`
	// \[appidx\] application index.
	ID uint64 `json:"id"`
	// Stores the global information associated with an application.
	Params *ApplicationParams `json:"params"`
}

// Stores local state associated with an application.
type ApplicationLocalState struct {
	// Round when account closed out of the application.
	ClosedOutAtRound *uint64 `json:"closedOutAtRound"`
	// Whether or not the application local state is currently deleted from its account.
	Deleted bool `json:"deleted"`
	// The application which this local state is for.
	ID uint64 `json:"id"`
	// Represents a key-value store for use in an application.
	KeyValue []TealKeyValue `json:"keyValue"`
	// Round when the account opted into the application.
	OptedInAtRound *uint64 `json:"optedInAtRound"`
	// Specifies maximums on the number of each type that may be stored.
	Schema *ApplicationStateSchema `json:"schema"`
}

// Stores the global information associated with an application.
type ApplicationParams struct {
	// \[approv\] approval program.
	ApprovalProgram []byte `json:"approvalProgram"`
	// \[clearp\] approval program.
	ClearStateProgram []byte `json:"clearStateProgram"`
	// The address that created this application. This is the address where the parameters and global state for this application can be found.
	Creator *string `json:"creator"`
	// \[epp\] the amount of extra program pages available to this app.
	ExtraProgramPages *uint64 `json:"extraProgramPages"`
	// Represents a key-value store for use in an application.
	GlobalState []TealKeyValue `json:"globalState"`
	// Specifies maximums on the number of each type that may be stored.
	GlobalStateSchema *ApplicationStateSchema `json:"globalStateSchema"`
	// Specifies maximums on the number of each type that may be stored.
	LocalStateSchema *ApplicationStateSchema `json:"localStateSchema"`
}

type ApplicationResponse struct {
	// Application index and its parameters
	Application *Application `json:"application"`
	// Round at which the results were computed.
	CurrentRound uint64 `json:"currentRound"`
}

// Specifies maximums on the number of each type that may be stored.
type ApplicationStateSchema struct {
	// \[nbs\] num of byte slices.
	NumByteSlice uint64 `json:"numByteSlice"`
	// \[nui\] num of uints.
	NumUint uint64 `json:"numUint"`
}

type ApplicationsResponse struct {
	Applications []Application `json:"applications"`
	// Round at which the results were computed.
	CurrentRound uint64 `json:"currentRound"`
	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken *string `json:"nextToken"`
}

// Specifies both the unique identifier and the parameters for an asset
type Asset struct {
	// Round during which this asset was created.
	CreatedAtRound *uint64 `json:"createdAtRound"`
	// Whether or not this asset is currently deleted.
	Deleted bool `json:"deleted"`
	// Round during which this asset was destroyed.
	DestroyedAtRound *uint64 `json:"destroyedAtRound"`
	// unique asset identifier
	ID uint64 `json:"id"`
	// AssetParams specifies the parameters for an asset.
	//
	// \[apar\] when part of an AssetConfig transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetParams
	Params *AssetParams `json:"params"`
}

type AssetBalancesResponse struct {
	Balances []MiniAssetHolding `json:"balances"`
	// Round at which the results were computed.
	CurrentRound uint64 `json:"currentRound"`
	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken *string `json:"nextToken"`
}

// Describes an asset held by an account.
//
// Definition:
// data/basics/userBalance.go : AssetHolding
type AssetHolding struct {
	// \[a\] number of units held.
	Amount uint64 `json:"amount"`
	// Asset ID of the holding.
	ID uint64 `json:"id"`
	// Address that created this asset. This is the address where the parameters for this asset can be found, and also the address where unwanted asset units can be sent in the worst case.
	Creator string `json:"creator"`
	// Whether or not the asset holding is currently deleted from its account.
	Deleted bool `json:"deleted"`
	// \[f\] whether or not the holding is frozen.
	Frozen bool `json:"frozen"`
	// Round during which the account opted into this asset holding.
	OptedInAtRound *uint64 `json:"optedInAtRound"`
	// Round during which the account opted out of this asset holding.
	OptedOutAtRound *uint64 `json:"optedOutAtRound"`
}

// AssetParams specifies the parameters for an asset.
//
// \[apar\] when part of an AssetConfig transaction.
//
// Definition:
// data/transactions/asset.go : AssetParams
type AssetParams struct {
	// \[c\] Address of account used to clawback holdings of this asset.  If empty, clawback is not permitted.
	Clawback *string `json:"clawback"`
	// The address that created this asset. This is the address where the parameters for this asset can be found, and also the address where unwanted asset units can be sent in the worst case.
	Creator string `json:"creator"`
	// \[dc\] The number of digits to use after the decimal point when displaying this asset. If 0, the asset is not divisible. If 1, the base unit of the asset is in tenths. If 2, the base unit of the asset is in hundredths, and so on. This value must be between 0 and 19 (inclusive).
	Decimals uint64 `json:"decimals"`
	// \[df\] Whether holdings of this asset are frozen by default.
	DefaultFrozen *bool `json:"defaultFrozen"`
	// \[f\] Address of account used to freeze holdings of this asset.  If empty, freezing is not permitted.
	Freeze *string `json:"freeze"`
	// \[m\] Address of account used to manage the keys of this asset and to destroy it.
	Manager *string `json:"manager"`
	// \[am\] A commitment to some unspecified asset metadata. The format of this metadata is up to the application.
	MetadataHash []byte `json:"metadataHash"`
	// \[an\] Name of this asset, as supplied by the creator.
	Name *string `json:"name"`
	// \[r\] Address of account holding reserve (non-minted) units of this asset.
	Reserve *string `json:"reserve"`
	// \[t\] The total number of units of this asset.
	Total uint64 `json:"total"`
	// \[un\] Name of a unit of this asset, as supplied by the creator.
	UnitName *string `json:"unitName"`
	// \[au\] URL where more information about the asset can be retrieved.
	URL *string `json:"url"`
}

type AssetResponse struct {
	// Specifies both the unique identifier and the parameters for an asset
	Asset *Asset `json:"asset"`
	// Round at which the results were computed.
	CurrentRound uint64 `json:"currentRound"`
}

type AssetsResponse struct {
	Assets []Asset `json:"assets"`
	// Round at which the results were computed.
	CurrentRound uint64 `json:"currentRound"`
	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken *string `json:"nextToken"`
}

// Block information.
//
// Definition:
// data/bookkeeping/block.go : Block
type Block struct {
	// \[gh\] hash to which this block belongs.
	GenesisHash []byte `json:"genesisHash"`
	// \[gen\] ID to which this block belongs.
	GenesisID string `json:"genesisId"`
	// \[prev\] Previous block hash.
	PreviousBlockHash []byte `json:"previousBlockHash"`
	// Fields relating to rewards,
	Rewards *BlockRewards `json:"rewards"`
	// \[rnd\] Current round on which this block was appended to the chain.
	Round uint64 `json:"round"`
	// \[seed\] Sortition seed.
	Seed []byte `json:"seed"`
	// \[ts\] Block creation timestamp in seconds since eposh
	Timestamp uint64 `json:"timestamp"`
	// \[txns\] list of transactions corresponding to a given round.
	Transactions []Transaction `json:"transactions"`
	// \[txn\] TransactionsRoot authenticates the set of transactions appearing in the block. More specifically, it's the root of a merkle tree whose leaves are the block's Txids, in lexicographic order. For the empty block, it's 0. Note that the TxnRoot does not authenticate the signatures on the transactions, only the transactions themselves. Two blocks with the same transactions but in a different order and with different signatures will have the same TxnRoot.
	TransactionsRoot []byte `json:"transactionsRoot"`
	// \[tc\] TxnCounter counts the number of transactions committed in the ledger, from the time at which support for this feature was introduced.
	//
	// Specifically, TxnCounter is the number of the next transaction that will be committed after this block.  It is 0 when no transactions have ever been committed (since TxnCounter started being supported).
	TxnCounter *uint64 `json:"txnCounter"`
	// Fields relating to a protocol upgrade.
	UpgradeState *BlockUpgradeState `json:"upgradeState"`
	// Fields relating to voting for a protocol upgrade.
	UpgradeVote *BlockUpgradeVote `json:"upgradeVote"`
}

// Fields relating to rewards,
type BlockRewards struct {
	// \[fees\] accepts transaction fees, it can only spend to the incentive pool.
	FeeSink string `json:"feeSink"`
	// \[rwcalr\] number of leftover MicroAlgos after the distribution of rewards-rate MicroAlgos for every reward unit in the next round.
	RewardsCalculationRound uint64 `json:"rewardsCalculationRound"`
	// \[earn\] How many rewards, in MicroAlgos, have been distributed to each RewardUnit of MicroAlgos since genesis.
	RewardsLevel uint64 `json:"rewardsLevel"`
	// \[rwd\] accepts periodic injections from the fee-sink and continually redistributes them as rewards.
	RewardsPool string `json:"rewardsPool"`
	// \[rate\] Number of new MicroAlgos added to the participation stake from rewards at the next round.
	RewardsRate uint64 `json:"rewardsRate"`
	// \[frac\] Number of leftover MicroAlgos after the distribution of RewardsRate/rewardUnits MicroAlgos for every reward unit in the next round.
	RewardsResidue uint64 `json:"rewardsResidue"`
}

// Fields relating to a protocol upgrade.
type BlockUpgradeState struct {
	// \[proto\] The current protocol version.
	CurrentProtocol string `json:"currentProtocol"`
	// \[nextproto\] The next proposed protocol version.
	NextProtocol *string `json:"nextProtocol"`
	// \[nextyes\] Number of blocks which approved the protocol upgrade.
	NextProtocolApprovals *uint64 `json:"nextProtocolApprovals"`
	// \[nextswitch\] Round on which the protocol upgrade will take effect.
	NextProtocolSwitchOn *uint64 `json:"nextProtocolSwitchOn"`
	// \[nextbefore\] Deadline round for this protocol upgrade (No votes will be consider after this round).
	NextProtocolVoteBefore *uint64 `json:"nextProtocolVoteBefore"`
}

// Fields relating to voting for a protocol upgrade.
type BlockUpgradeVote struct {
	// \[upgradeyes\] Indicates a yes vote for the current proposal.
	UpgradeApprove *bool `json:"upgradeApprove"`
	// \[upgradedelay\] Indicates the time between acceptance and execution.
	UpgradeDelay *uint64 `json:"upgradeDelay"`
	// \[upgradeprop\] Indicates a proposed upgrade.
	UpgradePropose *string `json:"upgradePropose"`
}

// Represents a TEAL value delta.
type EvalDelta struct {
	// \[at\] delta action.
	Action uint64 `json:"action"`
	// \[bs\] bytes value.
	Bytes []byte `json:"bytes"`
	// \[ui\] uint value.
	Uint *uint64 `json:"uint"`
}

// Key-value pairs for StateDelta.
type EvalDeltaKeyValue struct {
	Key []byte `json:"key"`
	// Represents a TEAL value delta.
	Value *EvalDelta `json:"value"`
}

// A health check response.
type HealthCheck struct {
	Data        map[string]interface{} `json:"data"`
	DbAvailable bool                   `json:"dbAvailable"`
	Errors      []string               `json:"errors"`
	IsMigrating bool                   `json:"isMigrating"`
	Message     string                 `json:"message"`
	Round       uint64                 `json:"round"`
}

// A simplified version of AssetHolding
type MiniAssetHolding struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
	// Whether or not this asset holding is currently deleted from its account.
	Deleted bool `json:"deleted"`
	Frozen  bool `json:"frozen"`
	// Round during which the account opted into the asset.
	OptedInAtRound *uint64 `json:"optedInAtRound"`
	// Round during which the account opted out of the asset.
	OptedOutAtRound *uint64 `json:"optedOutAtRound"`
}

// Represents a \[apls\] local-state or \[apgs\] global-state schema. These schemas determine how much storage may be used in a local-state or global-state for an application. The more space used, the larger minimum balance must be maintained in the account holding the data.
type StateSchema struct {
	// Maximum number of TEAL byte slices that may be stored in the key/value store.
	NumByteSlice uint64 `json:"numByteSlice"`
	// Maximum number of TEAL uints that may be stored in the key/value store.
	NumUint uint64 `json:"numUint"`
}

// Represents a key-value pair in an application store.
type TealKeyValue struct {
	Key []byte `json:"key"`
	// Represents a TEAL value.
	Value *TealValue `json:"value"`
}

// Represents a TEAL value.
type TealValue struct {
	// \[tb\] bytes value.
	Bytes []byte `json:"bytes"`
	// \[tt\] value type.
	Type uint64 `json:"type"`
	// \[ui\] uint value.
	Uint uint64 `json:"uint"`
}

// Contains all fields common to all transactions and serves as an envelope to all transactions type.
//
// Definition:
// data/transactions/signedtxn.go : SignedTxn
// data/transactions/transaction.go : Transaction
type Transaction struct {
	// Fields for application transactions.
	//
	// Definition:
	// data/transactions/application.go : ApplicationCallTxnFields
	ApplicationTransaction *TransactionApplication `json:"applicationTransaction"`
	// Fields for asset allocation, re-configuration, and destruction.
	//
	//
	// A zero value for asset-id indicates asset creation.
	// A zero value for the params indicates asset destruction.
	//
	// Definition:
	// data/transactions/asset.go : AssetConfigTxnFields
	AssetConfigTransaction *TransactionAssetConfig `json:"assetConfigTransaction"`
	// Fields for an asset freeze transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetFreezeTxnFields
	AssetFreezeTransaction *TransactionAssetFreeze `json:"assetFreezeTransaction"`
	// Fields for an asset transfer transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetTransferTxnFields
	AssetTransferTransaction *TransactionAssetTransfer `json:"assetTransferTransaction"`
	// \[sgnr\] this is included with signed transactions when the signing address does not equal the sender. The backend can use this to ensure that auth addr is equal to the accounts auth addr.
	AuthAddr *string `json:"authAddr"`
	// \[rc\] rewards applied to close-remainder-to account.
	CloseRewards *uint64 `json:"closeRewards"`
	// \[ca\] closing amount for transaction.
	ClosingAmount *uint64 `json:"closingAmount"`
	// Round when the transaction was confirmed.
	ConfirmedRound *uint64 `json:"confirmedRound"`
	// Specifies an application index (ID) if an application was created with this transaction.
	CreatedApplicationID *uint64 `json:"createdApplicationId"`
	// Specifies an asset index (ID) if an asset was created with this transaction.
	CreatedAssetID *uint64 `json:"createdAssetId"`
	// \[fee\] Transaction fee.
	Fee uint64 `json:"fee"`
	// \[fv\] First valid round for this transaction.
	FirstValid uint64 `json:"firstValid"`
	// \[gh\] Hash of genesis block.
	GenesisHash []byte `json:"genesisHash"`
	// \[gen\] genesis block ID.
	GenesisID *string `json:"genesisId"`
	// Application state delta.
	GlobalStateDelta []EvalDeltaKeyValue `json:"globalStateDelta"`
	// \[grp\] Base64 encoded byte array of a sha512/256 digest. When present indicates that this transaction is part of a transaction group and the value is the sha512/256 hash of the transactions in that group.
	Group []byte `json:"group"`
	// Transaction ID
	ID string `json:"id"`
	// Offset into the round where this transaction was confirmed.
	IntraRoundOffset *uint64 `json:"intraRoundOffset"`
	// Fields for a keyreg transaction.
	//
	// Definition:
	// data/transactions/keyreg.go : KeyregTxnFields
	KeyregTransaction *TransactionKeyreg `json:"keyregTransaction"`
	// \[lv\] Last valid round for this transaction.
	LastValid uint64 `json:"lastValid"`
	// \[lx\] Base64 encoded 32-byte array. Lease enforces mutual exclusion of transactions.  If this field is nonzero, then once the transaction is confirmed, it acquires the lease identified by the (Sender, Lease) pair of the transaction until the LastValid round passes.  While this transaction possesses the lease, no other transaction specifying this lease can be confirmed.
	Lease []byte `json:"lease"`
	// \[ld\] Local state key/value changes for the application being executed by this transaction.
	LocalStateDelta []AccountStateDelta `json:"localStateDelta"`
	// \[note\] Free form data.
	Note []byte `json:"note"`
	// Fields for a payment transaction.
	//
	// Definition:
	// data/transactions/payment.go : PaymentTxnFields
	PaymentTransaction *TransactionPayment `json:"paymentTransaction"`
	// \[rr\] rewards applied to receiver account.
	ReceiverRewards *uint64 `json:"receiverRewards"`
	// \[rekey\] when included in a valid transaction, the accounts auth addr will be updated with this value and future signatures must be signed with the key represented by this address.
	RekeyTo *string `json:"rekeyTo"`
	// Time when the block this transaction is in was confirmed.
	RoundTime *uint64 `json:"roundTime"`
	// \[snd\] Sender's address.
	Sender string `json:"sender"`
	// \[rs\] rewards applied to sender account.
	SenderRewards *uint64 `json:"senderRewards"`
	// Validation signature associated with some data. Only one of the signatures should be provided.
	Signature *TransactionSignature `json:"signature"`
	// \[type\] Indicates what type of transaction this is. Different types have different fields.
	//
	// Valid types, and where their fields are stored:
	// * \[pay\] payment-transaction
	// * \[keyreg\] keyreg-transaction
	// * \[acfg\] asset-config-transaction
	// * \[axfer\] asset-transfer-transaction
	// * \[afrz\] asset-freeze-transaction
	// * \[appl\] application-transaction
	TxType TxType `json:"txType"`
}

// Fields for application transactions.
//
// Definition:
// data/transactions/application.go : ApplicationCallTxnFields
type TransactionApplication struct {
	// \[apat\] List of accounts in addition to the sender that may be accessed from the application's approval-program and clear-state-program.
	Accounts []string `json:"accounts"`
	// \[apaa\] transaction specific arguments accessed from the application's approval-program and clear-state-program.
	ApplicationArgs [][]byte `json:"applicationArgs"`
	// \[apid\] ID of the application being configured or empty if creating.
	ApplicationID uint64 `json:"applicationId"`
	// \[apap\] Logic executed for every application transaction, except when on-completion is set to "clear". It can read and write global state for the application, as well as account-specific local state. Approval programs may reject the transaction.
	ApprovalProgram []byte `json:"approvalProgram"`
	// \[apsu\] Logic executed for application transactions with on-completion set to "clear". It can read and write global state for the application, as well as account-specific local state. Clear state programs cannot reject the transaction.
	ClearStateProgram []byte `json:"clearStateProgram"`
	// \[epp\] specifies the additional app program len requested in pages.
	ExtraProgramPages *uint64 `json:"extraProgramPages"`
	// \[apfa\] Lists the applications in addition to the application-id whose global states may be accessed by this application's approval-program and clear-state-program. The access is read-only.
	ForeignApps []uint64 `json:"foreignApps"`
	// \[apas\] lists the assets whose parameters may be accessed by this application's ApprovalProgram and ClearStateProgram. The access is read-only.
	ForeignAssets []uint64 `json:"foreignAssets"`
	// Represents a \[apls\] local-state or \[apgs\] global-state schema. These schemas determine how much storage may be used in a local-state or global-state for an application. The more space used, the larger minimum balance must be maintained in the account holding the data.
	GlobalStateSchema *StateSchema `json:"globalStateSchema"`
	// Represents a \[apls\] local-state or \[apgs\] global-state schema. These schemas determine how much storage may be used in a local-state or global-state for an application. The more space used, the larger minimum balance must be maintained in the account holding the data.
	LocalStateSchema *StateSchema `json:"localStateSchema"`
	// \[apan\] defines the what additional actions occur with the transaction.
	//
	// Valid types:
	// * noop
	// * optin
	// * closeout
	// * clear
	// * update
	// * update
	// * delete
	OnCompletion OnCompletion `json:"onCompletion"`
}

// Fields for asset allocation, re-configuration, and destruction.
//
//
// A zero value for asset-id indicates asset creation.
// A zero value for the params indicates asset destruction.
//
// Definition:
// data/transactions/asset.go : AssetConfigTxnFields
type TransactionAssetConfig struct {
	// \[xaid\] ID of the asset being configured or empty if creating.
	AssetID *uint64 `json:"assetId"`
	// AssetParams specifies the parameters for an asset.
	//
	// \[apar\] when part of an AssetConfig transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetParams
	Params *AssetParams `json:"params"`
}

// Fields for an asset freeze transaction.
//
// Definition:
// data/transactions/asset.go : AssetFreezeTxnFields
type TransactionAssetFreeze struct {
	// \[fadd\] Address of the account whose asset is being frozen or thawed.
	Address string `json:"address"`
	// \[faid\] ID of the asset being frozen or thawed.
	AssetID uint64 `json:"assetId"`
	// \[afrz\] The new freeze status.
	NewFreezeStatus bool `json:"newFreezeStatus"`
}

// Fields for an asset transfer transaction.
//
// Definition:
// data/transactions/asset.go : AssetTransferTxnFields
type TransactionAssetTransfer struct {
	// \[aamt\] Amount of asset to transfer. A zero amount transferred to self allocates that asset in the account's Assets map.
	Amount uint64 `json:"amount"`
	// \[xaid\] ID of the asset being transferred.
	AssetID uint64 `json:"assetId"`
	// Number of assets transfered to the close-to account as part of the transaction.
	CloseAmount *uint64 `json:"closeAmount"`
	// \[aclose\] Indicates that the asset should be removed from the account's Assets map, and specifies where the remaining asset holdings should be transferred.  It's always valid to transfer remaining asset holdings to the creator account.
	CloseTo *string `json:"closeTo"`
	// \[arcv\] Recipient address of the transfer.
	Receiver string `json:"receiver"`
	// \[asnd\] The effective sender during a clawback transactions. If this is not a zero value, the real transaction sender must be the Clawback address from the AssetParams.
	Sender *string `json:"sender"`
}

// Fields for a keyreg transaction.
//
// Definition:
// data/transactions/keyreg.go : KeyregTxnFields
type TransactionKeyreg struct {
	// \[nonpart\] Mark the account as participating or non-participating.
	NonParticipation bool `json:"nonParticipation"`
	// \[selkey\] Public key used with the Verified Random Function (VRF) result during committee selection.
	SelectionParticipationKey []byte `json:"selectionParticipationKey"`
	// \[votefst\] First round this participation key is valid.
	VoteFirstValid *uint64 `json:"voteFirstValid"`
	// \[votekd\] Number of subkeys in each batch of participation keys.
	VoteKeyDilution *uint64 `json:"voteKeyDilution"`
	// \[votelst\] Last round this participation key is valid.
	VoteLastValid *uint64 `json:"voteLastValid"`
	// \[votekey\] Participation public key used in key registration transactions.
	VoteParticipationKey []byte `json:"voteParticipationKey"`
}

// Fields for a payment transaction.
//
// Definition:
// data/transactions/payment.go : PaymentTxnFields
type TransactionPayment struct {
	// \[amt\] number of MicroAlgos intended to be transferred.
	Amount uint64 `json:"amount"`
	// Number of MicroAlgos that were sent to the close-remainder-to address when closing the sender account.
	CloseAmount *uint64 `json:"closeAmount"`
	// \[close\] when set, indicates that the sending account should be closed and all remaining funds be transferred to this address.
	CloseRemainderTo *string `json:"closeRemainderTo"`
	// \[rcv\] receiver's address.
	Receiver string `json:"receiver"`
}

type TransactionResponse struct {
	// Round at which the results were computed.
	CurrentRound uint64 `json:"currentRound"`
	// Contains all fields common to all transactions and serves as an envelope to all transactions type.
	//
	// Definition:
	// data/transactions/signedtxn.go : SignedTxn
	// data/transactions/transaction.go : Transaction
	Transaction *Transaction `json:"transaction"`
}

// Validation signature associated with some data. Only one of the signatures should be provided.
type TransactionSignature struct {
	// \[lsig\] Programatic transaction signature.
	//
	// Definition:
	// data/transactions/logicsig.go
	Logicsig *TransactionSignatureLogicsig `json:"logicsig"`
	// \[msig\] structure holding multiple subsignatures.
	//
	// Definition:
	// crypto/multisig.go : MultisigSig
	Multisig *TransactionSignatureMultisig `json:"multisig"`
	// \[sig\] Standard ed25519 signature.
	Sig []byte `json:"sig"`
}

// \[lsig\] Programatic transaction signature.
//
// Definition:
// data/transactions/logicsig.go
type TransactionSignatureLogicsig struct {
	// \[arg\] Logic arguments, base64 encoded.
	Args [][]byte `json:"args"`
	// \[l\] Program signed by a signature or multi signature, or hashed to be the address of ana ccount. Base64 encoded TEAL program.
	Logic []byte `json:"logic"`
	// \[msig\] structure holding multiple subsignatures.
	//
	// Definition:
	// crypto/multisig.go : MultisigSig
	MultisigSignature *TransactionSignatureMultisig `json:"multisigSignature"`
	// \[sig\] ed25519 signature.
	Signature []byte `json:"signature"`
}

// \[msig\] structure holding multiple subsignatures.
//
// Definition:
// crypto/multisig.go : MultisigSig
type TransactionSignatureMultisig struct {
	// \[subsig\] holds pairs of public key and signatures.
	Subsignature []TransactionSignatureMultisigSubsignature `json:"subsignature"`
	// \[thr\]
	Threshold *uint64 `json:"threshold"`
	// \[v\]
	Version *uint64 `json:"version"`
}

type TransactionSignatureMultisigSubsignature struct {
	// \[pk\]
	PublicKey []byte `json:"publicKey"`
	// \[s\]
	Signature []byte `json:"signature"`
}

type TransactionsResponse struct {
	// Round at which the results were computed.
	CurrentRound uint64 `json:"currentRound"`
	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken    *string       `json:"nextToken"`
	Transactions []Transaction `json:"transactions"`
}

type AccountStatus string

const (
	AccountStatusOffline          AccountStatus = "OFFLINE"
	AccountStatusOnline           AccountStatus = "ONLINE"
	AccountStatusNotParticipating AccountStatus = "NOT_PARTICIPATING"
)

var AllAccountStatus = []AccountStatus{
	AccountStatusOffline,
	AccountStatusOnline,
	AccountStatusNotParticipating,
}

func (e AccountStatus) IsValid() bool {
	switch e {
	case AccountStatusOffline, AccountStatusOnline, AccountStatusNotParticipating:
		return true
	}
	return false
}

func (e AccountStatus) String() string {
	return string(e)
}

func (e *AccountStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AccountStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AccountStatus", str)
	}
	return nil
}

func (e AccountStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type AddressRole string

const (
	AddressRoleSender       AddressRole = "SENDER"
	AddressRoleReceiver     AddressRole = "RECEIVER"
	AddressRoleFreezeTarget AddressRole = "FREEZE_TARGET"
)

var AllAddressRole = []AddressRole{
	AddressRoleSender,
	AddressRoleReceiver,
	AddressRoleFreezeTarget,
}

func (e AddressRole) IsValid() bool {
	switch e {
	case AddressRoleSender, AddressRoleReceiver, AddressRoleFreezeTarget:
		return true
	}
	return false
}

func (e AddressRole) String() string {
	return string(e)
}

func (e *AddressRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AddressRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AddressRole", str)
	}
	return nil
}

func (e AddressRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type OnCompletion string

const (
	OnCompletionNoop     OnCompletion = "NOOP"
	OnCompletionOptin    OnCompletion = "OPTIN"
	OnCompletionCloseout OnCompletion = "CLOSEOUT"
	OnCompletionClear    OnCompletion = "CLEAR"
	OnCompletionUpdate   OnCompletion = "UPDATE"
	OnCompletionDelete   OnCompletion = "DELETE"
)

var AllOnCompletion = []OnCompletion{
	OnCompletionNoop,
	OnCompletionOptin,
	OnCompletionCloseout,
	OnCompletionClear,
	OnCompletionUpdate,
	OnCompletionDelete,
}

func (e OnCompletion) IsValid() bool {
	switch e {
	case OnCompletionNoop, OnCompletionOptin, OnCompletionCloseout, OnCompletionClear, OnCompletionUpdate, OnCompletionDelete:
		return true
	}
	return false
}

func (e OnCompletion) String() string {
	return string(e)
}

func (e *OnCompletion) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = OnCompletion(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid OnCompletion", str)
	}
	return nil
}

func (e OnCompletion) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SigType string

const (
	SigTypeSig  SigType = "SIG"
	SigTypeMsig SigType = "MSIG"
	SigTypeLsig SigType = "LSIG"
)

var AllSigType = []SigType{
	SigTypeSig,
	SigTypeMsig,
	SigTypeLsig,
}

func (e SigType) IsValid() bool {
	switch e {
	case SigTypeSig, SigTypeMsig, SigTypeLsig:
		return true
	}
	return false
}

func (e SigType) String() string {
	return string(e)
}

func (e *SigType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SigType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SigType", str)
	}
	return nil
}

func (e SigType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TxType string

const (
	TxTypePay    TxType = "PAY"
	TxTypeKeyreg TxType = "KEYREG"
	TxTypeAcfg   TxType = "ACFG"
	TxTypeAxfer  TxType = "AXFER"
	TxTypeAfrz   TxType = "AFRZ"
	TxTypeAppl   TxType = "APPL"
)

var AllTxType = []TxType{
	TxTypePay,
	TxTypeKeyreg,
	TxTypeAcfg,
	TxTypeAxfer,
	TxTypeAfrz,
	TxTypeAppl,
}

func (e TxType) IsValid() bool {
	switch e {
	case TxTypePay, TxTypeKeyreg, TxTypeAcfg, TxTypeAxfer, TxTypeAfrz, TxTypeAppl:
		return true
	}
	return false
}

func (e TxType) String() string {
	return string(e)
}

func (e *TxType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TxType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TxType", str)
	}
	return nil
}

func (e TxType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
