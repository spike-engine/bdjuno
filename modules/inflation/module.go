package inflation

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spike-engine/juno/modules"

	"github.com/spike-engine/bdjuno/v3/database"
	inflationsource "github.com/spike-engine/bdjuno/v3/modules/inflation/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represent database/inflation module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source inflationsource.Source
}

// NewModule returns a new Module instance
func NewModule(source inflationsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "inflation"
}
