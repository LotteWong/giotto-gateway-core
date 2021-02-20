package utils

import (
	"crypto/sha256"
	"fmt"
	"hash"
)

func GenSaltPwd(salt, password string) string {
	var encrypter hash.Hash

	encrypter = sha256.New()
	encrypter.Write([]byte(password))
	cipherBeforeConcat := fmt.Sprintf("%x", encrypter.Sum(nil))

	encrypter = sha256.New()
	encrypter.Write([]byte(cipherBeforeConcat + salt))
	cipherAfterConcat := fmt.Sprintf("%x", encrypter.Sum(nil))

	return cipherAfterConcat
}
