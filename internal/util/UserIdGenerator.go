package util

import (
	"crypto/sha256"

	"github.com/google/uuid"
)

func GenerateUUIDFromString(id string) string {
	hash := sha256.New()
	hash.Write([]byte(id))
	return uuid.NewSHA1(uuid.NameSpaceOID, hash.Sum(nil)).String()
}