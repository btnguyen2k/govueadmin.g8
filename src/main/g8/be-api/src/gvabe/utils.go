package gvabe

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"strings"
)

const (
	systemGroupId = "administrator"

	systemAdminUsername = "admin"
	systemAdminName     = "Adam Local"
)

func encryptPassword(username, rawPassword string) string {
	saltAndPwd := username + "." + rawPassword
	out := sha1.Sum([]byte(saltAndPwd))
	return strings.ToLower(hex.EncodeToString(out[:]))
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

/*
randomString generates a random string with specified length.
*/
func randomString(l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
