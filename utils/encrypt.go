package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"

	"github.com/LotteWong/giotto-gateway-core/constants"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
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

func EncodeJwt(claims jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(constants.JwtSignKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func DecodeJwt(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.JwtSignKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("token is not jwt.StandardClaims")
	}

	return claims, nil
}
