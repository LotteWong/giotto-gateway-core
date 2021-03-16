package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
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

func MD5(text string) (string, error) {
	h := md5.New()
	if _, err := io.WriteString(h, text); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
