package urlshort

import (
	"crypto/sha256"
	"encoding/base64"
	"regexp"
)

func GetBaseEnc(url string) string {
	// We can use any hash function, sha256 is good vs collisiions
	sum := sha256.Sum256([]byte(url))
	// Use base64 for converting [0-9][a-f] to [0-9][a-z][A-Z]
	// Of course, base58/62 may be used. But base64 is std package
	encoded := base64.StdEncoding.EncodeToString(sum[:])
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	encoded = re.ReplaceAllString(encoded, "_")
	return encoded
}
