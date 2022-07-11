// dht22.go provides data structures and methods
// for control to dht22 sensor.
package dht

import (
	"encoding/binary"
	"strconv"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

// TryRead() cannot be sure to get sensor values.
// In many cases, dht22 sensor cannot be readable
// due to uncertain reason.
// So, many other libraries, in other languages,
// try to get for many times.
func (d *dht22Sensor) TryRead() (float32, float32, error) {
	return d.read()
}

// ReadRetry() provides more stable interface to get the values
// of sensor using automatically retrying.
// The recommended maxRetry's value is 15.
// Each trial be delayed in 2 seconds to improve probabilities.
// And each trial also opens and closes every time.
func (d *dht22Sensor) ReadRetry(maxRetry int) (float32, float32, error) {
	temperature, humidity, err := d.readRetry(maxRetry, 2)
	if err != nil {
		return 0, 0, err
	}
	return temperature, humidity, nil
}

// Internal data structures.
// DHT22 Sensor only needs pin number option.
// kidoman/embd library uses a string parameter to control GPIO
// For example,
// pin, err := embd.NewDigitalPin("GPIO_4")
// Therefore, eventually the pin number must be tranformed
// to string.
// So the pin number is a configuration not to change.
type dht22Config struct {
	Pin int
}

type dht22Sensor struct {
	pin embd.DigitalPin
	c   dht22Config
}

// These values are based on Adafruit python library.
const (
	DHT_PULSES   = 41
	DHT_MAXCOUNT = 32000
)

func (d *dht22Sensor) readRetry(retryCnt, delay int) (float32, float32, error) {
	var temperature, humidity float32
	var err error
	for i := 0; i < retryCnt; i++ {
		temperature, humidity, err = d.read()
		if err == nil {
			return temperature, humidity, nil
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
	return 0, 0, err
}

// TODO: refactor code
func (d *dht22Sensor) read() (float32, float32, error) {
	err := d.open()
	if err != nil {
		return 0, 0, err
	}
	defer d.close()

	pulseCounts := make([]int, DHT_PULSES*2)

	for {
		v, err := d.pin.Read()
		if err != nil {
			return 0, 0, err
		}
		if v != 1 {
			break
		}
	}

	for i := 0; i < DHT_PULSES*2; i += 2 {
		for {
			v, err := d.pin.Read()
			if err != nil {
				return 0, 0, err
			}
			if v != 0 {
				break
			}

			pulseCounts[i]++
			if pulseCounts[i] >= DHT_MAXCOUNT {
				return 0, 0, TimeoutError
			}
		}
		for {
			v, err := d.pin.Read()
			if err != nil {
				return 0, 0, err
			}
			if v != 1 {
				break
			}
			pulseCounts[i+1]++
			if pulseCounts[i+1] >= DHT_MAXCOUNT {
				return 0, 0, TimeoutError
			}
		}
	}
	threshold := 0
	for i := 2; i < DHT_PULSES*2; i += 2 {
		threshold += pulseCounts[i]
	}
	threshold /= DHT_PULSES - 1

	bytes := make([]uint8, 5)

	for i := 3; i < DHT_PULSES*2; i += 2 {
		idx := (i - 3) / 16
		bytes[idx] <<= 1
		if pulseCounts[i] >= threshold {
			bytes[idx] |= 1
		}
	}
	err = d.checksum(bytes)
	if err != nil {
		return 0, 0, err
	}

	temperature :=
		float32(int16(binary.BigEndian.Uint16(bytes[0:2]))) / 10
	humidity :=
		float32(int16(binary.BigEndian.Uint16(bytes[2:4]))) / 10

	if !ValidateHumidity(humidity) {
		return 0, 0, HumidityError
	}

	// datasheet operating range
	if !ValidateTemperature(temperature) {
		return 0, 0, TemperatureError
	}

	return temperature, humidity, nil
}

func (d *dht22Sensor) open() error {
	err := embd.InitGPIO()
	if err != nil {
		return err
	}

	gpioStr := "GPIO_" + strconv.Itoa(d.c.Pin)
	d.pin, err = embd.NewDigitalPin(gpioStr)
	if err != nil {
		return err
	}

	err = d.pin.SetDirection(embd.Out)
	if err != nil {
		return err
	}

	err = d.pin.Write(embd.High)
	if err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)

	err = d.pin.Write(embd.Low)
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Millisecond)

	err = d.pin.Write(embd.High)
	if err != nil {
		return err
	}
	time.Sleep(20 * time.Microsecond)

	err = d.pin.SetDirection(embd.In)
	if err != nil {
		return err
	}

	return nil
}

func (d *dht22Sensor) close() error {
	err := d.pin.Close()
	if err != nil {
		return err
	}

	err = embd.CloseGPIO()
	if err != nil {
		return err
	}

	return nil
}

func (d *dht22Sensor) checksum(bytes []uint8) error {
	sum := uint8(0)

	for i := 0; i < 4; i++ {
		sum += bytes[i]
	}

	if sum != bytes[4] {
		return ChecksumError
	}

	return nil
}

func newDefaultDHT22Sensor(pin int) *dht22Sensor {
	return newDHT22Sensor(dht22Config{
		Pin: pin,
	})
}

func newDHT22Sensor(c dht22Config) *dht22Sensor {
	return &dht22Sensor{
		pin: nil,
		c:   c,
	}
}
