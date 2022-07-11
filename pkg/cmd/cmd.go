package cmd

import (
	"fmt"

	"github.com/minuk-dev/rpi-humidity/pkg/dht"
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
	sensor := dht.New(dht.DHTConfig{
		Pin:  4,
		Type: dht.DHT22,
	})
	for {
		temperature, humidity, err := sensor.ReadRetry(15)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(temperature, humidity)
	}
}
