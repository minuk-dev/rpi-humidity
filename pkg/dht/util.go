package dht

func isValidHumidity(humidity float32) bool {
	return humidity >= 0 && humidity <= 100
}

func isValidTemperature(temperature float32) bool {
	return temperature >= -40 && temperature <= 80
}

func convertBytesToUint16(bytes []byte, isBigEndian bool) int16 {
	if isBigEndian {
		return int16(bytes[0])<<8 | int16(bytes[1])
	}
	return int16(bytes[1])<<8 | int16(bytes[0])
}
