package commands

import (
	"github.com/spf13/cobra"
)

func paramsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Manipulate parameters in the Parameter Registries",
	}
	cmd.AddCommand(
		paramsSettlementCmd(),
	)
	return cmd
}
