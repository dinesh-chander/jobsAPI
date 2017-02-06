package miscellaneous

import (
	"crypto/sha1"
	"encoding/hex"
)

func GenerateSHAChecksum(value string) (checksum string) {
	newSHA := sha1.New()

	newSHA.Write([]byte(value))

	checksum = hex.EncodeToString(newSHA.Sum(nil))

	return
}
