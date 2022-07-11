package cli

import (
	"github.com/spf13/cobra"
)

func RunNoErrorOutput(cmd *cobra.Command) error {
	err := cmd.Execute()
	if err != nil {
		return err
	}

	return nil
}
