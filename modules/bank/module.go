package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/spike-engine/bdjuno/v3/database"
	"github.com/spike-engine/bdjuno/v3/modules/bank/source"

	junomessages "github.com/spike-engine/juno/modules/messages"

	"github.com/spike-engine/juno/modules"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represents the x/bank module
type Module struct {
	cdc codec.Codec
	db  *database.Db

	messageParser junomessages.MessageAddressesParser
	keeper        source.Source
}

// NewModule returns a new Module instance
func NewModule(
	messageParser junomessages.MessageAddressesParser, keeper source.Source, cdc codec.Codec, db *database.Db,
) *Module {
	return &Module{
		cdc:           cdc,
		db:            db,
		messageParser: messageParser,
		keeper:        keeper,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "bank"
}
