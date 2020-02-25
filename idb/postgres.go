// Copyright (C) 2019-2020 Algorand, Inc.
// This file is part of the Algorand Indexer
//
// Algorand Indexer is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// Algorand Indexer is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Algorand Indexer.  If not, see <https://www.gnu.org/licenses/>.

// You can build without postgres by `go build --tags nopostgres` but it's on by default
// +build !nopostgres

package idb

// import text to contstant setup_postgres_sql
//go:generate go run ../cmd/texttosource/main.go idb setup_postgres.sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/json"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	atypes "github.com/algorand/go-algorand-sdk/types"
	_ "github.com/lib/pq"

	"github.com/algorand/indexer/types"
)

func OpenPostgres(connection string) (pdb *PostgresIndexerDb, err error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	pdb = &PostgresIndexerDb{
		db:         db,
		protoCache: make(map[string]types.ConsensusParams, 20),
	}
	err = pdb.init()
	return
}

type PostgresIndexerDb struct {
	db *sql.DB
	tx *sql.Tx

	protoCache map[string]types.ConsensusParams
}

func (db *PostgresIndexerDb) init() (err error) {
	_, err = db.db.Exec(setup_postgres_sql)
	return
}

func (db *PostgresIndexerDb) AlreadyImported(path string) (imported bool, err error) {
	row := db.db.QueryRow(`SELECT COUNT(path) FROM imported WHERE path = $1`, path)
	numpath := 0
	err = row.Scan(&numpath)
	return numpath == 1, err
}

func (db *PostgresIndexerDb) MarkImported(path string) (err error) {
	_, err = db.db.Exec(`INSERT INTO imported (path) VALUES ($1)`, path)
	return err
}

func (db *PostgresIndexerDb) StartBlock() (err error) {
	db.tx, err = db.db.BeginTx(context.Background(), nil)
	return
}

func (db *PostgresIndexerDb) AddTransaction(round uint64, intra int, txtypeenum int, assetid uint64, txnbytes []byte, txn types.SignedTxnInBlock, participation [][]byte) error {
	// TODO: set txn_participation
	var err error
	txid := crypto.TransactionIDString(txn.Txn)
	_, err = db.tx.Exec(`INSERT INTO txn (round, intra, typeenum, asset, txid, txnbytes, txn) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING`, round, intra, txtypeenum, assetid, txid[:], txnbytes, string(json.Encode(txn)))
	if err != nil {
		return err
	}
	stmt, err := db.tx.Prepare(`INSERT INTO txn_participation (addr, round, intra) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`)
	if err != nil {
		return err
	}
	for _, paddr := range participation {
		_, err = stmt.Exec(paddr, round, intra)
		if err != nil {
			return err
		}
	}
	return err
}
func (db *PostgresIndexerDb) CommitBlock(round uint64, timestamp int64, rewardslevel uint64, headerbytes []byte) error {
	var err error
	_, err = db.tx.Exec(`INSERT INTO block_header (round, realtime, rewardslevel, header) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`, round, time.Unix(timestamp, 0), rewardslevel, headerbytes)
	if err != nil {
		return err
	}
	err = db.tx.Commit()
	db.tx = nil
	return err
}

func (db *PostgresIndexerDb) GetBlockHeader(round uint64) (block types.Block, err error) {
	row := db.db.QueryRow(`SELECT header FROM block_header WHERE round = $1`, round)
	var blockbytes []byte
	err = row.Scan(&blockbytes)
	if err != nil {
		return
	}
	err = msgpack.Decode(blockbytes, &block)
	return
}

// GetAsset return AssetParams about an asset
func (db *PostgresIndexerDb) GetAsset(assetid uint64) (asset types.AssetParams, err error) {
	row := db.db.QueryRow(`SELECT params FROM asset WHERE index = $1`, assetid)
	var assetjson string
	err = row.Scan(&assetjson)
	if err != nil {
		return
	}
	err = json.Decode([]byte(assetjson), &asset)
	return
}

// GetDefaultFrozen get {assetid:default frozen, ...} for all assets
func (db *PostgresIndexerDb) GetDefaultFrozen() (defaultFrozen map[uint64]bool, err error) {
	rows, err := db.db.Query(`SELECT index, params -> 'df' FROM asset`)
	if err != nil {
		return
	}
	defaultFrozen = make(map[uint64]bool)
	for rows.Next() {
		var assetid uint64
		var frozen bool
		err = rows.Scan(&assetid, &frozen)
		if err != nil {
			return
		}
		defaultFrozen[assetid] = frozen
	}
	return
}

func (db *PostgresIndexerDb) LoadGenesis(genesis types.Genesis) (err error) {
	tx, err := db.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback() // ignored if .Commit() first

	setAccount, err := tx.Prepare(`INSERT INTO account (addr, microalgos, rewardsbase, account_data) VALUES ($1, $2, 0, $3)`)
	if err != nil {
		return
	}
	defer setAccount.Close()

	total := uint64(0)
	for ai, alloc := range genesis.Allocation {
		addr, err := atypes.DecodeAddress(alloc.Address)
		if len(alloc.State.AssetParams) > 0 || len(alloc.State.Assets) > 0 {
			return fmt.Errorf("genesis account[%d] has unhandled asset", ai)
		}
		_, err = setAccount.Exec(addr[:], alloc.State.MicroAlgos, string(json.Encode(alloc.State)))
		total += uint64(alloc.State.MicroAlgos)
		if err != nil {
			return fmt.Errorf("error setting genesis account[%d], %v", ai, err)
		}
	}
	err = tx.Commit()
	fmt.Printf("genesis %d accounts %d microalgos, %v\n", len(genesis.Allocation), total, err)
	return err

}

func (db *PostgresIndexerDb) SetProto(version string, proto types.ConsensusParams) (err error) {
	pj := json.Encode(proto)
	if err != nil {
		return err
	}
	_, err = db.db.Exec(`INSERT INTO protocol (version, proto) VALUES ($1, $2) ON CONFLICT (version) DO UPDATE SET proto = EXCLUDED.proto`, version, pj)
	return err
}

func (db *PostgresIndexerDb) GetProto(version string) (proto types.ConsensusParams, err error) {
	proto, hit := db.protoCache[version]
	if hit {
		return
	}
	row := db.db.QueryRow(`SELECT proto FROM protocol WHERE version = $1`, version)
	var protostr string
	err = row.Scan(&protostr)
	if err != nil {
		return
	}
	err = json.Decode([]byte(protostr), &proto)
	if err == nil {
		db.protoCache[version] = proto
	}
	return
}

func (db *PostgresIndexerDb) GetMetastate(key string) (jsonStrValue string, err error) {
	row := db.db.QueryRow(`SELECT v FROM metastate WHERE k = $1`, key)
	err = row.Scan(&jsonStrValue)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (db *PostgresIndexerDb) SetMetastate(key, jsonStrValue string) (err error) {
	_, err = db.db.Exec(`INSERT INTO metastate (k, v) VALUES ($1, $2) ON CONFLICT (k) DO UPDATE SET v = EXCLUDED.v`, key, jsonStrValue)
	return
}

func (db *PostgresIndexerDb) GetMaxRound() (round uint64, err error) {
	row := db.db.QueryRow(`SELECT max(round) FROM block_header`)
	err = row.Scan(&round)
	return
}

func (db *PostgresIndexerDb) yieldTxnsThread(ctx context.Context, rows *sql.Rows, results chan<- TxnRow) {
	for rows.Next() {
		var round uint64
		var intra int
		var txnbytes []byte
		err := rows.Scan(&round, &intra, &txnbytes)
		var row TxnRow
		if err != nil {
			row.Error = err
		} else {
			row.Round = round
			row.Intra = intra
			row.TxnBytes = txnbytes
		}
		select {
		case <-ctx.Done():
			break
		case results <- row:
			if err != nil {
				break
			}
		}
	}
	close(results)
}

func (db *PostgresIndexerDb) YieldTxns(ctx context.Context, prevRound int64) <-chan TxnRow {
	results := make(chan TxnRow, 1)
	rows, err := db.db.QueryContext(ctx, `SELECT round, intra, txnbytes FROM txn WHERE round > $1 ORDER BY round, intra`, prevRound)
	if err != nil {
		results <- TxnRow{Error: err}
		close(results)
		return results
	}
	go db.yieldTxnsThread(ctx, rows, results)
	return results
}

func (db *PostgresIndexerDb) CommitRoundAccounting(updates RoundUpdates, round, rewardsBase uint64) (err error) {
	any := false
	tx, err := db.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback() // ignored if .Commit() first

	if len(updates.AlgoUpdates) > 0 {
		any = true
		// account_data json is only used on account creation, otherwise the account data jsonb field is updated from the delta
		setalgo, err := tx.Prepare(`INSERT INTO account (addr, microalgos, rewardsbase) VALUES ($1, $2, $3) ON CONFLICT (addr) DO UPDATE SET microalgos = account.microalgos + EXCLUDED.microalgos, rewardsbase = EXCLUDED.rewardsbase`)
		if err != nil {
			return fmt.Errorf("prepare update algo, %v", err)
		}
		defer setalgo.Close()
		for addr, delta := range updates.AlgoUpdates {
			_, err = setalgo.Exec(addr[:], delta, rewardsBase)
			if err != nil {
				return fmt.Errorf("update algo, %v", err)
			}
		}
	}
	if len(updates.AcfgUpdates) > 0 {
		any = true
		setacfg, err := tx.Prepare(`INSERT INTO asset (index, creator_addr, params) VALUES ($1, $2, $3) ON CONFLICT (index) DO UPDATE SET params = EXCLUDED.params`)
		if err != nil {
			return fmt.Errorf("prepare set asset, %v", err)
		}
		defer setacfg.Close()
		for _, au := range updates.AcfgUpdates {
			_, err = setacfg.Exec(au.AssetId, au.Creator[:], string(json.Encode(au.Params)))
			if err != nil {
				return fmt.Errorf("update asset, %v", err)
			}
		}
	}
	if len(updates.AssetUpdates) > 0 {
		any = true
		seta, err := tx.Prepare(`INSERT INTO account_asset (addr, assetid, amount, frozen) VALUES ($1, $2, $3, $4) ON CONFLICT (addr, assetid) DO UPDATE SET amount = account_asset.amount + EXCLUDED.amount`)
		if err != nil {
			return fmt.Errorf("prepare set account_asset, %v", err)
		}
		defer seta.Close()
		for _, au := range updates.AssetUpdates {
			_, err = seta.Exec(au.Addr[:], au.AssetId, au.Delta, au.DefaultFrozen)
			if err != nil {
				return fmt.Errorf("update account asset, %v", err)
			}
		}
	}
	if len(updates.FreezeUpdates) > 0 {
		any = true
		fr, err := tx.Prepare(`INSERT INTO account_asset (addr, assetid, amount, frozen) VALUES ($1, $2, 0, $3) ON CONFLICT (addr, assetid) DO UPDATE SET frozen = EXCLUDED.frozen`)
		if err != nil {
			return fmt.Errorf("prepare asset freeze, %v", err)
		}
		defer fr.Close()
		for _, fs := range updates.FreezeUpdates {
			_, err = fr.Exec(fs.Addr[:], fs.AssetId, fs.Frozen)
			if err != nil {
				return fmt.Errorf("update asset freeze, %v", err)
			}
		}
	}
	if len(updates.AssetCloses) > 0 {
		any = true
		acs, err := tx.Prepare(`INSERT INTO account_asset (addr, assetid, amount, frozen)
SELECT $1, $2, x.amount, $3 FROM account_asset x WHERE x.addr = $4
ON CONFLICT (addr, assetid) DO UPDATE SET amount = account_asset.amount + EXCLUDED.amount`)
		if err != nil {
			return fmt.Errorf("prepare asset close1, %v", err)
		}
		defer acs.Close()
		acd, err := tx.Prepare(`DELETE FROM account_asset WHERE addr = $1`)
		if err != nil {
			return fmt.Errorf("prepare asset close2, %v", err)
		}
		defer acd.Close()
		for _, ac := range updates.AssetCloses {
			_, err = acs.Exec(ac.CloseTo[:], ac.AssetId, ac.DefaultFrozen, ac.Sender[:])
			if err != nil {
				return fmt.Errorf("asset close send, %v", err)
			}
			_, err = acd.Exec(ac.Sender[:])
			if err != nil {
				return fmt.Errorf("asset close del, %v", err)
			}
		}
	}
	if len(updates.AssetDestroys) > 0 {
		any = true
		// Note! leaves `asset` row present for historical reference, but deletes all holdings from all accounts
		ads, err := tx.Prepare(`DELETE FROM account_asset WHERE assetid = $1`)
		if err != nil {
			return fmt.Errorf("prepare asset destroy, %v", err)
		}
		defer ads.Close()
		for _, assetId := range updates.AssetDestroys {
			ads.Exec(assetId)
			if err != nil {
				return fmt.Errorf("asset destroy, %v", err)
			}
		}
	}
	if !any {
		fmt.Printf("empty round %d\n", round)
	}
	var istate ImportState
	staterow := tx.QueryRow(`SELECT v FROM metastate WHERE k = 'state'`)
	var stateJsonStr string
	err = staterow.Scan(&stateJsonStr)
	if err == sql.ErrNoRows {
		// ok
	} else if err != nil {
		return
	} else {
		err = json.Decode([]byte(stateJsonStr), &istate)
		if err != nil {
			return
		}
	}
	istate.AccountRound = int64(round)
	sjs := string(json.Encode(istate))
	_, err = tx.Exec(`INSERT INTO metastate (k, v) VALUES ('state', $1) ON CONFLICT (k) DO UPDATE SET v = EXCLUDED.v`, sjs)
	if err != nil {
		return
	}
	return tx.Commit()
}

func (db *PostgresIndexerDb) GetBlock(round uint64) (block types.Block, err error) {
	row := db.db.QueryRow(`SELECT header FROM block_header WHERE round = $1`, round)
	var blockheaderbytes []byte
	err = row.Scan(&blockheaderbytes)
	if err != nil {
		return
	}
	err = msgpack.Decode(blockheaderbytes, &block)
	return
}

func (db *PostgresIndexerDb) Transactions(ctx context.Context, tf TransactionFilter) <-chan TxnRow {
	const maxWhereParts = 14
	whereParts := make([]string, 0, maxWhereParts)
	whereArgs := make([]interface{}, 0, maxWhereParts)
	joinParticipation := false
	joinHeader := false
	partNumber := 1
	if tf.Address != nil {
		whereParts = append(whereParts, fmt.Sprintf("p.addr = $%d", partNumber))
		whereArgs = append(whereArgs, tf.Address)
		partNumber++
		joinParticipation = true
	}
	if tf.MinRound != 0 {
		whereParts = append(whereParts, fmt.Sprintf("t.round >= $%d", partNumber))
		whereArgs = append(whereArgs, tf.MinRound)
		partNumber++
	}
	if tf.MaxRound != 0 {
		whereParts = append(whereParts, fmt.Sprintf("t.round <= $%d", partNumber))
		whereArgs = append(whereArgs, tf.MaxRound)
		partNumber++
	}
	if !tf.BeforeTime.IsZero() {
		whereParts = append(whereParts, fmt.Sprintf("h.realtime < $%d", partNumber))
		whereArgs = append(whereArgs, tf.BeforeTime)
		partNumber++
		joinHeader = true
	}
	if !tf.AfterTime.IsZero() {
		whereParts = append(whereParts, fmt.Sprintf("h.realtime > $%d", partNumber))
		whereArgs = append(whereArgs, tf.AfterTime)
		partNumber++
		joinHeader = true
	}
	if tf.AssetId != 0 {
		whereParts = append(whereParts, fmt.Sprintf("t.asset = $%d", partNumber))
		whereArgs = append(whereArgs, tf.AssetId)
		partNumber++
	}
	if tf.TypeEnum != 0 {
		whereParts = append(whereParts, fmt.Sprintf("t.typeenum = $%d", partNumber))
		whereArgs = append(whereArgs, tf.TypeEnum)
		partNumber++
	}
	if tf.Txid != nil {
		whereParts = append(whereParts, fmt.Sprintf("t.txid = $%d", partNumber))
		whereArgs = append(whereArgs, tf.Txid)
		partNumber++
	}
	if tf.Round >= 0 {
		whereParts = append(whereParts, fmt.Sprintf("t.round = $%d", partNumber))
		whereArgs = append(whereArgs, tf.Round)
		partNumber++
	}
	if tf.Offset >= 0 {
		whereParts = append(whereParts, fmt.Sprintf("t.intra = $%d", partNumber))
		whereArgs = append(whereArgs, tf.Offset)
		partNumber++
	}
	if len(tf.SigType) != 0 {
		whereParts = append(whereParts, fmt.Sprintf("t.txn -> $%d IS NOT NULL", partNumber))
		whereArgs = append(whereArgs, tf.SigType)
		partNumber++
	}
	query := "SELECT t.round, t.intra, t.txnbytes FROM txn t"
	whereStr := strings.Join(whereParts, " AND ")
	if joinParticipation {
		query += " JOIN txn_participation p ON t.round = p.round AND t.intra = p.intra"
	}
	if joinHeader {
		query += " JOIN block_header h ON t.round = h.round"
	}
	query += " WHERE " + whereStr
	if joinParticipation {
		// this should match the index on txn_particpation
		query += " ORDER BY p.addr, p.round DESC, p.intra DESC"
	} else {
		// this should explicitly match the primary key on txn (round,intra)
		query += " ORDER BY t.round, t.intra"
	}
	if tf.Limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", tf.Limit)
	}
	out := make(chan TxnRow, 1)
	rows, err := db.db.QueryContext(ctx, query, whereArgs...)
	if err != nil {
		err = fmt.Errorf("txn query %#v err %v", query, err)
		out <- TxnRow{Error: err}
		close(out)
		return out
	}
	go db.yieldTxnsThread(ctx, rows, out)
	return out
}

const maxAccountsLimit = 1000

func (db *PostgresIndexerDb) yieldAccountsThread(ctx context.Context, opts AccountQueryOptions, rows *sql.Rows, tx *sql.Tx, blockheader types.Block, out chan<- AccountRow) {
	defer tx.Rollback()
	for rows.Next() {
		var addr []byte
		var microalgos uint64
		var rewardsbase uint64
		var accountDataJsonStr *string

		// these are bytes of json serialization
		var holdingAssetid []byte
		var holdingAmount []byte
		var holdingFrozen []byte

		// these are bytes of json serialization
		var assetParamsIds []byte
		var assetParamsStr []byte

		var err error

		if opts.IncludeAssetHoldings {
			if opts.IncludeAssetParams {
				err = rows.Scan(
					&addr, &microalgos, &rewardsbase, &accountDataJsonStr,
					&holdingAssetid, &holdingAmount, &holdingFrozen,
					&assetParamsIds, &assetParamsStr,
				)
			} else {
				err = rows.Scan(
					&addr, &microalgos, &rewardsbase, &accountDataJsonStr,
					&holdingAssetid, &holdingAmount, &holdingFrozen,
				)
			}
		} else if opts.IncludeAssetParams {
			err = rows.Scan(
				&addr, &microalgos, &rewardsbase, &accountDataJsonStr,
				&assetParamsIds, &assetParamsStr,
			)
		} else {
			err = rows.Scan(&addr, &microalgos, &rewardsbase, &accountDataJsonStr)
		}
		if err != nil {
			out <- AccountRow{Error: err}
			break
		}

		var account models.Account
		var aaddr atypes.Address
		copy(aaddr[:], addr)
		account.Address = aaddr.String()
		account.Round = uint64(blockheader.Round)
		account.AmountWithoutPendingRewards = microalgos

		// TODO: pending rewards calculation doesn't belong in database layer (this is just the most covenient place which has all the data)
		proto, err := db.GetProto(string(blockheader.CurrentProtocol))
		rewardsUnits := uint64(0)
		if proto.RewardUnit != 0 {
			rewardsUnits = microalgos / proto.RewardUnit
		}
		rewardsDelta := blockheader.RewardsLevel - rewardsbase
		account.PendingRewards = rewardsUnits * rewardsDelta
		account.Amount = microalgos + account.PendingRewards
		// not implemented: account.Rewards sum of all rewards ever

		const nullarraystr = "[null]"

		if len(holdingAssetid) > 0 && string(holdingAssetid) != nullarraystr {
			var haids []uint64
			err = json.Decode(holdingAssetid, &haids)
			if err != nil {
				out <- AccountRow{Error: err}
				break
			}
			var hamounts []uint64
			err = json.Decode(holdingAmount, &hamounts)
			if err != nil {
				out <- AccountRow{Error: err}
				break
			}
			var hfrozen []bool
			err = json.Decode(holdingFrozen, &hfrozen)
			if err != nil {
				out <- AccountRow{Error: err}
				break
			}
			account.Assets = make(map[uint64]models.AssetHolding, len(haids))
			for i, assetid := range haids {
				account.Assets[assetid] = models.AssetHolding{Amount: hamounts[i], Frozen: hfrozen[i]}
			}
		}
		if len(assetParamsIds) > 0 && string(assetParamsIds) != nullarraystr {
			var assetids []uint64
			err = json.Decode(assetParamsIds, &assetids)
			if err != nil {
				out <- AccountRow{Error: err}
				break
			}
			var assetParams []types.AssetParams
			err = json.Decode(assetParamsStr, &assetParams)
			if err != nil {
				out <- AccountRow{Error: err}
				break
			}
			account.AssetParams = make(map[uint64]models.AssetParams, len(assetids))
			for i, assetid := range assetids {
				ap := assetParams[i]
				account.AssetParams[assetid] = models.AssetParams{
					Creator:       account.Address,
					Total:         ap.Total,
					Decimals:      ap.Decimals,
					DefaultFrozen: ap.DefaultFrozen,
					UnitName:      ap.UnitName,
					AssetName:     ap.AssetName,
					URL:           ap.URL,
					MetadataHash:  ap.MetadataHash[:],
					ManagerAddr:   addrStr(ap.Manager[:]),
					ReserveAddr:   addrStr(ap.Reserve[:]),
					FreezeAddr:    addrStr(ap.Freeze[:]),
					ClawbackAddr:  addrStr(ap.Clawback[:]),
				}

			}
		}
		out <- AccountRow{Account: account}
	}
	close(out)
}

func addrStr(addr []byte) string {
	if len(addr) == 0 {
		return ""
	}
	allZero := true
	for _, bv := range addr {
		if bv != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		return ""
	}
	var aa atypes.Address
	copy(aa[:], addr)
	return aa.String()
}

func (db *PostgresIndexerDb) GetAccounts(ctx context.Context, opts AccountQueryOptions) <-chan AccountRow {
	out := make(chan AccountRow, 1)

	// Begin transaction so we get everything at one consistent point in time and round of accounting.
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("account tx err %v", err)
		out <- AccountRow{Error: err}
		close(out)
		return out
	}

	// Get round number through which accounting has been updated
	row := tx.QueryRow(`SELECT (v -> 'account_round')::bigint as account_round FROM metastate WHERE k = 'state'`)
	var accountRound uint64
	err = row.Scan(&accountRound)
	if err != nil {
		err = fmt.Errorf("account_round err %v", err)
		out <- AccountRow{Error: err}
		close(out)
		tx.Rollback()
		return out
	}

	// Get block header for that round so we know protocol and rewards info
	row = tx.QueryRow(`SELECT header FROM block_header WHERE round = $1`, accountRound)
	var headerbytes []byte
	err = row.Scan(&headerbytes)
	if err != nil {
		err = fmt.Errorf("account round header %d err %v", accountRound, err)
		out <- AccountRow{Error: err}
		close(out)
		tx.Rollback()
		return out
	}
	var blockheader types.Block
	err = msgpack.Decode(headerbytes, &blockheader)
	if err != nil {
		err = fmt.Errorf("account round header %d err %v", accountRound, err)
		out <- AccountRow{Error: err}
		close(out)
		tx.Rollback()
		return out
	}

	// Construct query for fetching accounts...
	query := `SELECT a.addr, a.microalgos, a.rewardsbase, a.account_data`
	if opts.IncludeAssetHoldings {
		query += `, json_agg(aa.assetid) as haid, json_agg(aa.amount) as hamt, json_agg(aa.frozen) as hf`
	}
	if opts.IncludeAssetParams {
		query += `, json_agg(ap.index) as paid, json_agg(ap.params) as pp`
	}
	query += ` FROM account a`
	if opts.IncludeAssetHoldings {
		query += ` LEFT JOIN account_asset aa ON a.addr = aa.addr`
	}
	if opts.IncludeAssetParams {
		query += ` LEFT JOIN asset ap ON a.addr = ap.creator_addr`
	}
	const maxWhereParts = 14
	whereParts := make([]string, 0, maxWhereParts)
	whereArgs := make([]interface{}, 0, maxWhereParts)
	partNumber := 1
	if len(opts.GreaterThanAddress) > 0 {
		whereParts = append(whereParts, fmt.Sprintf("a.addr > $%d", partNumber))
		whereArgs = append(whereArgs, opts.GreaterThanAddress)
		partNumber++
	}
	if len(opts.EqualToAddress) > 0 {
		whereParts = append(whereParts, fmt.Sprintf("a.addr = $%d", partNumber))
		whereArgs = append(whereArgs, opts.EqualToAddress)
		partNumber++
	}
	if len(whereParts) > 0 {
		whereStr := strings.Join(whereParts, " AND ")
		query += " WHERE " + whereStr
	}
	if opts.IncludeAssetHoldings || opts.IncludeAssetParams {
		query += " GROUP BY 1,2,3,4"
	}
	if opts.Limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	rows, err := tx.Query(query, whereArgs...)
	if err != nil {
		err = fmt.Errorf("account query %#v err %v", query, err)
		out <- AccountRow{Error: err}
		close(out)
		tx.Rollback()
		return out
	}
	go db.yieldAccountsThread(ctx, opts, rows, tx, blockheader, out)
	return out
}

type postgresFactory struct {
}

func (df postgresFactory) Name() string {
	return "postgres"
}
func (df postgresFactory) Build(arg string) (IndexerDb, error) {
	return OpenPostgres(arg)
}

func init() {
	indexerFactories = append(indexerFactories, &postgresFactory{})
}

type ImportState struct {
	AccountRound int64 `codec:"account_round"`
}

func ParseImportState(js string) (istate ImportState, err error) {
	err = json.Decode([]byte(js), &istate)
	return
}