// Package dht provides simple dht(digital humidity and temperature sensor) interface.
//
// It can only supports dht22 sensor now.
package dht

import (
	"github.com/pkg/errors"
)

// Supported sensor types.
// They are candidates for the value of DHTConfig's Type.
const (
	DHT22 int = 0
)

// Generic DHT sensor errors.
// Errors returned by Read() and ReadRetry() can be tested
// against these errors using errors.Is.
var (
	ChecksumError       = errors.New("checksum error")
	HumidityError       = errors.New("humidity range error")
	TemperatureError    = errors.New("temperature range error")
	TooManyRequestError = errors.New("too many requests error")
	TimeoutError        = errors.New("timeout error")
)

// DHTSensor provides two interfaces to get temperature and humidity due to the probability of failure.
//
// Example Code
// sensor := dht.New(dht.DHTConfig{ Pin: 4, Type: dht.DHT22 })
// temperature, humidity, err := sensor.Read()
type DHTSensor interface {
	TryRead() (float32, float32, error)
	ReadRetry(maxRetryCnt int) (float32, float32, error)
}

// DHTConfig contains GPIO pin and sensor type for DHTSensor
type DHTConfig struct {
	Pin  int
	Type int
}

// New returns a new instance of DHTSensor from the given config.
func New(o DHTConfig) DHTSensor {
	if o.Type != DHT22 {
		return nil
	}
	return newDefaultDHT22Sensor(o.Pin)
}
