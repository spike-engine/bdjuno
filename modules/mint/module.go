package mint

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spike-engine/juno/modules"

	"github.com/spike-engine/bdjuno/v3/database"
	mintsource "github.com/spike-engine/bdjuno/v3/modules/mint/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/mint module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source mintsource.Source
}

// NewModule returns a new Module instance
func NewModule(source mintsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "mint"
}
