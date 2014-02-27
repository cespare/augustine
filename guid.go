package augustine

import (
	"crypto/rand"
)

const guidLen = 16

// guid generates a 128-bit random value
func guid() ([]byte, error) {
	id := make([]byte, guidLen)
	if _, err := rand.Read(id); err != nil {
		return nil, err
	}
	return id, nil
}
