package memory

import (
	"errors"
	"fmt"
	"urlshort/pkg/hashfunc"
)

// FromHash возвращает полный URL из мапы
func FromHash(hmap map[string]string, hash string) (string, error) {
	if hmap == nil {
		return "", errors.New("memory: FromHash(): hash table is nil")
	}
	if value, exists := hmap[hash]; exists {
		return value, nil
	} else {
		return "", fmt.Errorf(`no url for hash "%s" found`, hash)
	}
}

// ToHash принимает полный URL, вносит полный URL и сокращенный URL в мапу
func ToHash(hmap map[string]string, url string) (hashBaseEnc string, err error) {
	if hmap == nil {
		return "", errors.New("memory: ToHash(): hash table is nil")
	}
	hashBaseEnc = hashfunc.GetBaseEnc(url)
	for i := 10; i < 32; i++ {
		if value, exists := hmap[hashBaseEnc[i-10:i]]; exists {
			if value == url {
				/* This URL is already on hash rable */
				return hashBaseEnc[i-10 : i], errors.New("url is already on hash table")
			}
		} else {
			hmap[hashBaseEnc[i-10:i]] = url
			return hashBaseEnc[i-10 : i], nil
		}
	}
	return "", errors.New("collision not resolved")
}
