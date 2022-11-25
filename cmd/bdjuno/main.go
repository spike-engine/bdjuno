package main

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spike-engine/juno/cmd"
	initcmd "github.com/spike-engine/juno/cmd/init"
	parsetypes "github.com/spike-engine/juno/cmd/parse/types"
	startcmd "github.com/spike-engine/juno/cmd/start"
	"github.com/spike-engine/juno/modules/messages"

	migratecmd "github.com/spike-engine/bdjuno/v3/cmd/migrate"
	parsecmd "github.com/spike-engine/bdjuno/v3/cmd/parse"

	"github.com/spike-engine/bdjuno/v3/types/config"

	"github.com/spike-engine/bdjuno/v3/database"
	"github.com/spike-engine/bdjuno/v3/modules"

	gaiaapp "github.com/cosmos/gaia/v7/app"
	evmosapp "github.com/evmos/evmos/v6/app"
)

func main() {
	initCfg := initcmd.NewConfig().
		WithConfigCreator(config.Creator)

	parseCfg := parsetypes.NewConfig().
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(config.MakeEncodingConfig(getBasicManagers())).
		WithRegistrar(modules.NewRegistrar(getAddressesParser()))

	cfg := cmd.NewConfig("bdjuno").
		WithInitConfig(initCfg).
		WithParseConfig(parseCfg)

	// Run the command
	rootCmd := cmd.RootCmd(cfg.GetName())

	rootCmd.AddCommand(
		cmd.VersionCmd(),
		initcmd.NewInitCmd(cfg.GetInitConfig()),
		parsecmd.NewParseCmd(cfg.GetParseConfig()),
		migratecmd.NewMigrateCmd(cfg.GetName(), cfg.GetParseConfig()),
		startcmd.NewStartCmd(cfg.GetParseConfig()),
	)

	executor := cmd.PrepareRootCmd(cfg.GetName(), rootCmd)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

// getBasicManagers returns the various basic managers that are used to register the encoding to
// support custom messages.
// This should be edited by custom implementations if needed.
func getBasicManagers() []module.BasicManager {
	return []module.BasicManager{
		gaiaapp.ModuleBasics,
		evmosapp.ModuleBasics,
	}
}

// getAddressesParser returns the messages parser that should be used to get the users involved in
// a specific message.
// This should be edited by custom implementations if needed.
func getAddressesParser() messages.MessageAddressesParser {
	return messages.JoinMessageParsers(
		messages.CosmosMessageAddressesParser,
	)
}
