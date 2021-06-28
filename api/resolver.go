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

	newBlockHeaderLock      sync.Mutex
	newBlockHeaderListeners map[chan<- *model.BlockHeader]bool
}

func NewResolver(si *ServerImplementation) *Resolver {
	resolver := Resolver{
		si:                      si,
		newBlockHeaderListeners: make(map[chan<- *model.BlockHeader]bool),
	}

	si.db.SetBlockCommitHook(resolver.onBlockCommit)

	return &resolver
}

func (r *Resolver) addBlockHeaderListener(ctx context.Context) <-chan *model.BlockHeader {
	events := make(chan *model.BlockHeader, 1)

	go func() {
		<-ctx.Done()
		r.newBlockHeaderLock.Lock()
		delete(r.newBlockHeaderListeners, events)
		r.newBlockHeaderLock.Unlock()
		close(events)
	}()

	r.newBlockHeaderLock.Lock()
	r.newBlockHeaderListeners[events] = true
	r.newBlockHeaderLock.Unlock()

	return events
}

func (r *Resolver) onBlockCommit(blockHeader types.BlockHeader) {
	r.newBlockHeaderLock.Lock()
	converted := helper.InternalBlockHeaderToModel(blockHeader)
	for listener := range r.newBlockHeaderListeners {
		select {
		case listener <- converted:
		default:
			// could not send message to client
			// remove them from map?
		}
	}
	r.newBlockHeaderLock.Unlock()
}
