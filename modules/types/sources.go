package types

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/simapp"
	inflationtypes "github.com/evmos/evmos/v6/x/inflation/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/spike-engine/juno/node/remote"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spike-engine/juno/node/local"

	nodeconfig "github.com/spike-engine/juno/node/config"

	banksource "github.com/spike-engine/bdjuno/v3/modules/bank/source"
	localbanksource "github.com/spike-engine/bdjuno/v3/modules/bank/source/local"
	remotebanksource "github.com/spike-engine/bdjuno/v3/modules/bank/source/remote"
	distrsource "github.com/spike-engine/bdjuno/v3/modules/distribution/source"
	localdistrsource "github.com/spike-engine/bdjuno/v3/modules/distribution/source/local"
	remotedistrsource "github.com/spike-engine/bdjuno/v3/modules/distribution/source/remote"
	govsource "github.com/spike-engine/bdjuno/v3/modules/gov/source"
	localgovsource "github.com/spike-engine/bdjuno/v3/modules/gov/source/local"
	remotegovsource "github.com/spike-engine/bdjuno/v3/modules/gov/source/remote"
	inflationsource "github.com/spike-engine/bdjuno/v3/modules/inflation/source"
	localinflationsource "github.com/spike-engine/bdjuno/v3/modules/inflation/source/local"
	remoteinflationsource "github.com/spike-engine/bdjuno/v3/modules/inflation/source/remote"
	mintsource "github.com/spike-engine/bdjuno/v3/modules/mint/source"
	localmintsource "github.com/spike-engine/bdjuno/v3/modules/mint/source/local"
	remotemintsource "github.com/spike-engine/bdjuno/v3/modules/mint/source/remote"
	slashingsource "github.com/spike-engine/bdjuno/v3/modules/slashing/source"
	localslashingsource "github.com/spike-engine/bdjuno/v3/modules/slashing/source/local"
	remoteslashingsource "github.com/spike-engine/bdjuno/v3/modules/slashing/source/remote"
	stakingsource "github.com/spike-engine/bdjuno/v3/modules/staking/source"
	localstakingsource "github.com/spike-engine/bdjuno/v3/modules/staking/source/local"
	remotestakingsource "github.com/spike-engine/bdjuno/v3/modules/staking/source/remote"

	evmosapp "github.com/evmos/evmos/v6/app"
)

type Sources struct {
	BankSource      banksource.Source
	DistrSource     distrsource.Source
	GovSource       govsource.Source
	InflationSource inflationsource.Source
	MintSource      mintsource.Source
	SlashingSource  slashingsource.Source
	StakingSource   stakingsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return buildLocalSources(cfg, encodingConfig)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, encodingConfig *params.EncodingConfig) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	app := simapp.NewSimApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{},
	)

	evmosApp := evmosapp.NewEvmos(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{},
	)

	sources := &Sources{
		BankSource:      localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		DistrSource:     localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistrKeeper)),
		GovSource:       localgovsource.NewSource(source, govtypes.QueryServer(app.GovKeeper)),
		InflationSource: localinflationsource.NewSource(source, inflationtypes.QueryServer(evmosApp.InflationKeeper)),
		MintSource:      localmintsource.NewSource(source, minttypes.QueryServer(app.MintKeeper)),
		SlashingSource:  localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		StakingSource:   localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: app.StakingKeeper}),
	}

	// Mount and initialize the stores
	err = source.MountKVStores(app, "keys")
	if err != nil {
		return nil, err
	}

	err = source.MountTransientStores(app, "tkeys")
	if err != nil {
		return nil, err
	}

	err = source.MountMemoryStores(app, "memKeys")
	if err != nil {
		return nil, err
	}

	err = source.InitStores()
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		BankSource:      remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:     remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:       remotegovsource.NewSource(source, govtypes.NewQueryClient(source.GrpcConn)),
		InflationSource: remoteinflationsource.NewSource(source, inflationtypes.NewQueryClient(source.GrpcConn)),
		MintSource:      remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource:  remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:   remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
