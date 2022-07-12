package db

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdbAPI "github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxDBConfig struct {
	ServerURL string
	AuthToken string
	Org       string
	Bucket    string
}

type influxDB struct {
	client   influxdb2.Client
	config   InfluxDBConfig
	writeAPI influxdbAPI.WriteAPI
}

func New(c InfluxDBConfig) *influxDB {
	return &influxDB{
		config: c,
	}
}

func (db *influxDB) Open() error {
	db.client = influxdb2.NewClient(db.config.ServerURL, db.config.AuthToken)
	_, err := db.client.Ready(context.TODO())
	if err != nil {
		return err
	}

	db.writeAPI = db.client.WriteAPI(db.config.Org, db.config.Bucket)
	return nil
}

func (db *influxDB) Write(m Measurement) error {
	p := influxdb2.NewPoint("home",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{
			"temperature": m.Temperature,
			"humidity":    m.Humidity,
		},
		time.Now())
	db.writeAPI.WritePoint(p)
	db.writeAPI.Flush()
	return nil
}

func (db *influxDB) Close() {
	db.client.Close()
}
