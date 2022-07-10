package dht

import "github.com/pkg/errors"

const (
	DHT22 int = 0
)

var (
	ChecksumError       = errors.New("checksum error")
	HumidityError       = errors.New("humidity range error")
	TemperatureError    = errors.New("temperature range error")
	TooManyRequestError = errors.New("too many requests error")
	TimeoutError        = errors.New("timeout error")
)

type DHTSensor interface {
	TryRead() (float32, float32, error)
	ReadRetry(maxRetryCnt int) (float32, float32, error)
}

type DHTOptions struct {
	Pin  int
	Type int
}

func New(o DHTOptions) DHTSensor {
	if o.Type != DHT22 {
		return nil
	}
	return newDefaultDHT22Sensor(o.Pin)
}
