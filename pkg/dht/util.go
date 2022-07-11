// util.go provides common functions.
package dht

// ValidateHumidity checks that the value of humidity are
// in normal range.
// Every sensor has own their capable range, but humidity
// is a kind of percentage.
func ValidateHumidity(humidity float32) bool {
	return humidity >= 0 && humidity <= 100
}

// ValidateTemperature checks that the value of temperature
// are in normal range.
// Every sensor has own their capable range.
// So, it can be absorbed into sensor's method.
func ValidateTemperature(temperature float32) bool {
	return temperature >= -40 && temperature <= 80
}
