package auth

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/spike-engine/bdjuno/v3/database"

	"github.com/spike-engine/juno/modules"
	"github.com/spike-engine/juno/modules/messages"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/auth module
type Module struct {
	cdc            codec.Codec
	db             *database.Db
	messagesParser messages.MessageAddressesParser
}

// NewModule builds a new Module instance
func NewModule(messagesParser messages.MessageAddressesParser, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		messagesParser: messagesParser,
		cdc:            cdc,
		db:             db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "auth"
}
