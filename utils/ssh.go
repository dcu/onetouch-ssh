package utils

import (
	"fmt"

	"crypto/md5"
	"golang.org/x/crypto/ssh"
)

// PublicKeyFingerprint Helper function to get the fingerprint of a SSH public key
func PublicKeyFingerprint(pubkey string) (string, string) {
	byteString := []byte(pubkey)

	pk, comment, _, _, err := ssh.ParseAuthorizedKey(byteString)
	if err != nil {
		return "", ""
	}

	hash := md5.Sum(pk.Marshal())
	fingerprint := Rfc4716Hex(hash[:])

	return fingerprint, comment
}

// Rfc4716Hex returns the data in hex format separated by ":"
func Rfc4716Hex(data []byte) string {
	var fingerprint string
	for i := 0; i < len(data); i++ {
		fingerprint = fmt.Sprintf("%s%0.2x", fingerprint, data[i])
		if i != len(data)-1 {
			fingerprint = fingerprint + ":"
		}
	}

	return fingerprint
}
