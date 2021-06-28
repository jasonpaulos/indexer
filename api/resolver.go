package api

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"context"
	"sync"

	"github.com/algorand/indexer/api/graph/helper"
	"github.com/algorand/indexer/api/graph/model"
	"github.com/algorand/indexer/types"
)

type Resolver struct {
	si *ServerImplementation

	newBlockLock      sync.Mutex
	newBlockListeners map[chan<- *model.Block]bool

	accountUpdateLock      sync.Mutex
	accountUpdateListeners map[string][]chan<- *model.AccountUpdateResponse
}

func NewResolver(si *ServerImplementation) *Resolver {
	resolver := Resolver{
		si:                     si,
		newBlockListeners:      make(map[chan<- *model.Block]bool),
		accountUpdateListeners: make(map[string][]chan<- *model.AccountUpdateResponse),
	}

	si.db.SetBlockCommitHook(resolver.onBlockCommit)

	return &resolver
}

func (r *Resolver) addBlockListener(ctx context.Context) <-chan *model.Block {
	events := make(chan *model.Block, 1)

	go func() {
		<-ctx.Done()
		r.newBlockLock.Lock()
		delete(r.newBlockListeners, events)
		r.newBlockLock.Unlock()
		close(events)
	}()

	r.newBlockLock.Lock()
	r.newBlockListeners[events] = true
	r.newBlockLock.Unlock()

	return events
}

func (r *Resolver) addAccountListener(ctx context.Context, address string) <-chan *model.AccountUpdateResponse {
	events := make(chan *model.AccountUpdateResponse, 1)

	go func() {
		<-ctx.Done()
		r.accountUpdateLock.Lock()

		listeners := r.accountUpdateListeners[address]
		index := -1
		for i, listener := range listeners {
			if listener == events {
				index = i
				break
			}
		}

		if index != -1 {
			listeners[index] = listeners[len(listeners)-1]
			listeners = listeners[:len(listeners)-1]
		}

		if len(listeners) == 0 {
			delete(r.accountUpdateListeners, address)
		} else {
			r.accountUpdateListeners[address] = listeners
		}

		r.accountUpdateLock.Unlock()
		close(events)
	}()

	r.accountUpdateLock.Lock()
	r.accountUpdateListeners[address] = append(r.accountUpdateListeners[address], events)
	r.accountUpdateLock.Unlock()

	return events
}

func (r *Resolver) onBlockCommit(header types.BlockHeader, txns []types.SignedTxnWithAD, createdPrimitives map[int]uint64, participants map[string][]int) {
	block := helper.InternalBlockHeaderAndTxnsToModel(header, txns, createdPrimitives)

	r.newBlockLock.Lock()
	for listener := range r.newBlockListeners {
		select {
		case listener <- block:
		default:
			// could not send message to client
		}
	}
	r.newBlockLock.Unlock()

	r.accountUpdateLock.Lock()
	for account, listeners := range r.accountUpdateListeners {
		updates := participants[account]
		if len(updates) == 0 {
			continue
		}

		txns := make([]model.Transaction, len(updates))
		for i, txnIndex := range updates {
			txns[i] = block.Transactions[txnIndex]
		}

		response := model.AccountUpdateResponse{
			Transactions: txns,
		}

		for _, listener := range listeners {
			select {
			case listener <- &response:
			default:
				// could not send message to client
			}
		}
	}
	r.accountUpdateLock.Unlock()
}
