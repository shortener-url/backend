package utils

import (
	"time"

	"gitlab.com/uchile1/helper/helperLog"
)

func Base62Encode(keyLen int) string {
	timestamp := time.Now().UnixNano()

	const base = 62
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!()*+-/<=>@[]^_{|}~"

	var encoded string
	for timestamp > 0 {
		remainder := timestamp % base
		encoded = string(chars[remainder]) + encoded
		timestamp = timestamp / base
	}
	aux := encoded
	if len(encoded) >= keyLen {
		encoded = encoded[len(encoded)-keyLen:] // Tomar los últimos 5 caracteres
	} else if len(encoded) < keyLen {
		// Rellenar con caracteres a la izquierda si la longitud es menor
		for len(encoded) < keyLen {
			encoded = "0" + encoded
		}
	}
	helperLog.Logger.Info().Msgf("URL: %s -> %s", aux, encoded)
	return encoded
}
