// Package handler defines the contract bot handlers implement to mount their
// routes on the shared telegofx router.
package handler

import "github.com/go-core-fx/telegofx"

// Handler registers bot routes on the shared telegofx router; each handler
// declares its own predicate chain and telegohandler dispatches
// first-match-wins over registration order.
type Handler interface {
	// Register mounts the handler's routes on router.
	Register(router *telegofx.Router)
}
