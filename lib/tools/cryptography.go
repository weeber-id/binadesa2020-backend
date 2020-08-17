package tools

import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeMD5 encrypt
func EncodeMD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte("testing"))
	return hex.EncodeToString(hasher.Sum(nil))
}
