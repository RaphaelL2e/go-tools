package main

import (
	"crypto/md5"
	"fmt"
)

var salt = "this is salt"

func Encryption(str string) string {
	data := []byte(str)

	md5str1 := fmt.Sprintf("%x", md5.Sum(data))

	return md5str1
}
func EncryptionSalt(str string) string {
	str = str + salt
	data := []byte(str)

	md5str1 := fmt.Sprintf("%x", md5.Sum(data))

	return md5str1
}

func SetSalt(nsalt string) {
	salt = nsalt
}
