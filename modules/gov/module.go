package gov

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/spike-engine/bdjuno/v3/database"

	govsource "github.com/spike-engine/bdjuno/v3/modules/gov/source"

	"github.com/spike-engine/juno/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent x/gov module
type Module struct {
	cdc             codec.Codec
	db              *database.Db
	source          govsource.Source
	authModule      AuthModule
	distrModule     DistrModule
	inflationModule InflationModule
	mintModule      MintModule
	slashingModule  SlashingModule
	stakingModule   StakingModule
}

// NewModule returns a new Module instance
func NewModule(
	source govsource.Source,
	authModule AuthModule,
	distrModule DistrModule,
	inflationModule InflationModule,
	mintModule MintModule,
	slashingModule SlashingModule,
	stakingModule StakingModule,
	cdc codec.Codec,
	db *database.Db,
) *Module {
	return &Module{
		cdc:             cdc,
		source:          source,
		authModule:      authModule,
		distrModule:     distrModule,
		inflationModule: inflationModule,
		mintModule:      mintModule,
		slashingModule:  slashingModule,
		stakingModule:   stakingModule,
		db:              db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "gov"
}
