package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type RPIHumidityOptions struct {
}

func NewDefaultRPIHumidityCommand() *cobra.Command {
	return NewDefaultRPIHumidityCommandWithArgs(RPIHumidityOptions{})
}

func NewDefaultRPIHumidityCommandWithArgs(o RPIHumidityOptions) *cobra.Command {
	cmd := NewRPIHumidityCommand(o)
	return cmd
}

func NewRPIHumidityCommand(o RPIHumidityOptions) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "rpi-humidity",
		Short: "rpi-humidity logger",
		Run:   runRPIHumidity,
	}
	return cmds
}

func runRPIHumidity(cmd *cobra.Command, args []string) {
	fmt.Println("Not Implemented")
}
