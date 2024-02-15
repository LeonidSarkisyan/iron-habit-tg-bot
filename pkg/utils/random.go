package utils

import (
	"crypto/rand"
	"encoding/binary"
)

func GenerateRandomSeed() (int64, error) {
	var seed int64
	err := binary.Read(rand.Reader, binary.BigEndian, &seed)
	if err != nil {
		return 0, err
	}
	return seed, nil
}
