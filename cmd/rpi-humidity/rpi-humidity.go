package main

import (
	"github.com/minuk-dev/rpi-humidity/pkg/cli"
	"github.com/minuk-dev/rpi-humidity/pkg/cmd"
	"github.com/minuk-dev/rpi-humidity/pkg/util"
)

func main() {
	command := cmd.NewDefaultRPIHumidityCommand()
	if err := cli.RunNoErrorOutput(command); err != nil {
		util.CheckErr(err)
	}
}
