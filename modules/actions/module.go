package actions

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/spike-engine/juno/modules"
	"github.com/spike-engine/juno/node"
	"github.com/spike-engine/juno/node/builder"
	nodeconfig "github.com/spike-engine/juno/node/config"
	"github.com/spike-engine/juno/types/config"

	modulestypes "github.com/spike-engine/bdjuno/v3/modules/types"
)

const (
	ModuleName = "actions"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
)

type Module struct {
	cfg     *Config
	node    node.Node
	sources *modulestypes.Sources
}

func NewModule(cfg config.Config, encodingConfig *params.EncodingConfig) *Module {
	bz, err := cfg.GetBytes()
	if err != nil {
		panic(err)
	}

	actionsCfg, err := ParseConfig(bz)
	if err != nil {
		panic(err)
	}

	nodeCfg := cfg.Node
	if actionsCfg.Node != nil {
		nodeCfg = nodeconfig.NewConfig(nodeconfig.TypeRemote, actionsCfg.Node)
	}

	// Build the node
	junoNode, err := builder.BuildNode(nodeCfg, encodingConfig)
	if err != nil {
		panic(err)
	}

	// Build the sources
	sources, err := modulestypes.BuildSources(nodeCfg, encodingConfig)
	if err != nil {
		panic(err)
	}

	return &Module{
		cfg:     actionsCfg,
		node:    junoNode,
		sources: sources,
	}
}

func (m *Module) Name() string {
	return ModuleName
}
