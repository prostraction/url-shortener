package hashfunc

import (
	"crypto/sha256"
	"encoding/base64"
	"regexp"
)

// GetBaseEnc принимает на вход URL и возвращает его хеш, приведенный к Base64
// Для Base64 строки спец. символы заменены на _
func GetBaseEnc(url string) string {
	// Мы можем использовать любую хеш функцию, sha256 хороша против коллизий
	sum := sha256.Sum256([]byte(url))
	// Используем base64 для конвертации [0-9][a-f] в [0-9][a-z][A-Z]
	// Можем использовать и base62, конечно, но base64 в std пакете
	encoded := base64.StdEncoding.EncodeToString(sum[:])
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	encoded = re.ReplaceAllString(encoded, "_")
	return encoded
}
