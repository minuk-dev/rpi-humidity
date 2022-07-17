package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/minuk-dev/rpi-humidity/pkg/db"
	"github.com/minuk-dev/rpi-humidity/pkg/dht"
	"github.com/spf13/cobra"
)

type RPIHumidityOptions struct {
	DBServer  string
	AuthToken string
	Org       string
	Bucket    string
}

func NewDefaultRPIHumidityCommand() *cobra.Command {
	return NewDefaultRPIHumidityCommandWithArgs(RPIHumidityOptions{})
}

func NewDefaultRPIHumidityCommandWithArgs(o RPIHumidityOptions) *cobra.Command {
	cmd := NewRPIHumidityCommand(o)
	return cmd
}

func NewRPIHumidityCommand(o RPIHumidityOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rpi-humidity",
		Short: "rpi-humidity logger",
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Validate()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = o.Run()
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
	cmd.Flags().StringVarP(&o.DBServer, "database-server", "d", "localhost:8086", "influxdb serverURL")
	cmd.Flags().StringVarP(&o.AuthToken, "database-token", "t", "", "influxdb auth token")
	cmd.Flags().StringVarP(&o.Org, "org", "o", "", "")
	cmd.Flags().StringVarP(&o.Bucket, "bucket", "b", "", "")
	return cmd
}

func (o *RPIHumidityOptions) Validate() error {
	return nil
}

func (o *RPIHumidityOptions) Run() error {
	client := db.New(db.InfluxDBConfig{
		ServerURL: o.DBServer,
		AuthToken: o.AuthToken,
		Org:       o.Org,
		Bucket:    o.Bucket,
	})

	err := client.Open()
	if err != nil {
		return err
	}
	defer client.Close()

	sensor := dht.New(dht.DHTConfig{
		Pin:  4,
		Type: dht.DHT22,
	})
	if sensor == nil {
		return errors.New("cannot initialize dht")
	}

	for {
		temperature, humidity, err := sensor.ReadRetry(15)
		if err != nil {
			fmt.Println(err)
		}
		err = client.Write(db.Measurement{
			Temperature: temperature,
			Humidity:    humidity,
		})
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(10)
	}

	return nil
}
