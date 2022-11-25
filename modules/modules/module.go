package modules

import (
	"github.com/spike-engine/juno/modules"
	"github.com/spike-engine/juno/types/config"

	"github.com/spike-engine/bdjuno/v3/database"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
)

type Module struct {
	cfg config.ChainConfig
	db  *database.Db
}

// NewModule returns a new Module instance
func NewModule(cfg config.ChainConfig, db *database.Db) *Module {
	return &Module{
		cfg: cfg,
		db:  db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "modules"
}
