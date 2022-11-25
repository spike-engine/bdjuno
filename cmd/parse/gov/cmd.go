package gov

import (
	"github.com/spf13/cobra"
	parsecmdtypes "github.com/spike-engine/juno/cmd/parse/types"
)

// NewGovCmd returns the Cobra command allowing to fix various things related to the x/gov module
func NewGovCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gov",
		Short: "Fix things related to the x/gov module",
	}

	cmd.AddCommand(
		proposalCmd(parseConfig),
	)

	return cmd
}
