package modules

import (
	"github.com/spike-engine/bdjuno/v3/modules/actions"
	"github.com/spike-engine/bdjuno/v3/modules/types"

	"github.com/spike-engine/juno/modules/pruning"
	"github.com/spike-engine/juno/modules/telemetry"

	"github.com/spike-engine/bdjuno/v3/modules/slashing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jmodules "github.com/spike-engine/juno/modules"
	"github.com/spike-engine/juno/modules/messages"
	"github.com/spike-engine/juno/modules/registrar"

	"github.com/spike-engine/bdjuno/v3/utils"

	"github.com/spike-engine/bdjuno/v3/database"
	"github.com/spike-engine/bdjuno/v3/modules/auth"
	"github.com/spike-engine/bdjuno/v3/modules/bank"
	"github.com/spike-engine/bdjuno/v3/modules/consensus"
	"github.com/spike-engine/bdjuno/v3/modules/distribution"
	"github.com/spike-engine/bdjuno/v3/modules/feegrant"
	"github.com/spike-engine/bdjuno/v3/modules/inflation"

	"github.com/spike-engine/bdjuno/v3/modules/gov"
	"github.com/spike-engine/bdjuno/v3/modules/mint"
	"github.com/spike-engine/bdjuno/v3/modules/modules"
	"github.com/spike-engine/bdjuno/v3/modules/pricefeed"
	"github.com/spike-engine/bdjuno/v3/modules/staking"
)

// UniqueAddressesParser returns a wrapper around the given parser that removes all duplicated addresses
func UniqueAddressesParser(parser messages.MessageAddressesParser) messages.MessageAddressesParser {
	return func(cdc codec.Codec, msg sdk.Msg) ([]string, error) {
		addresses, err := parser(cdc, msg)
		if err != nil {
			return nil, err
		}

		return utils.RemoveDuplicateValues(addresses), nil
	}
}

// --------------------------------------------------------------------------------------------------------------------

var (
	_ registrar.Registrar = &Registrar{}
)

// Registrar represents the modules.Registrar that allows to register all modules that are supported by BigDipper
type Registrar struct {
	parser messages.MessageAddressesParser
}

// NewRegistrar allows to build a new Registrar instance
func NewRegistrar(parser messages.MessageAddressesParser) *Registrar {
	return &Registrar{
		parser: UniqueAddressesParser(parser),
	}
}

// BuildModules implements modules.Registrar
func (r *Registrar) BuildModules(ctx registrar.Context) jmodules.Modules {
	cdc := ctx.EncodingConfig.Marshaler
	db := database.Cast(ctx.Database)

	sources, err := types.BuildSources(ctx.JunoConfig.Node, ctx.EncodingConfig)
	if err != nil {
		panic(err)
	}

	actionsModule := actions.NewModule(ctx.JunoConfig, ctx.EncodingConfig)
	authModule := auth.NewModule(r.parser, cdc, db)
	bankModule := bank.NewModule(r.parser, sources.BankSource, cdc, db)
	consensusModule := consensus.NewModule(db)
	distrModule := distribution.NewModule(sources.DistrSource, cdc, db)
	feegrantModule := feegrant.NewModule(cdc, db)
	inflationModule := inflation.NewModule(sources.InflationSource, cdc, db)
	mintModule := mint.NewModule(sources.MintSource, cdc, db)
	slashingModule := slashing.NewModule(sources.SlashingSource, cdc, db)
	stakingModule := staking.NewModule(sources.StakingSource, slashingModule, cdc, db)
	govModule := gov.NewModule(sources.GovSource, authModule, distrModule, inflationModule, mintModule, slashingModule, stakingModule, cdc, db)

	return []jmodules.Module{
		messages.NewModule(r.parser, cdc, ctx.Database),
		telemetry.NewModule(ctx.JunoConfig),
		pruning.NewModule(ctx.JunoConfig, db, ctx.Logger),

		actionsModule,
		authModule,
		bankModule,
		consensusModule,
		distrModule,
		feegrantModule,
		govModule,
		inflationModule,
		mintModule,
		modules.NewModule(ctx.JunoConfig.Chain, db),
		pricefeed.NewModule(ctx.JunoConfig, cdc, db),
		slashingModule,
		stakingModule,
	}
}
