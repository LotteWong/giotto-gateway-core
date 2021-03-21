package utils

import (
	"encoding/json"
	"strings"
)

func Obj2Json(obj interface{}) string {
	bts, _ := json.Marshal(obj)
	return string(bts)
}

func JoinSlash(a, b string) string {
	suffixSlash := strings.HasSuffix(a, "/")
	prefixSlash := strings.HasPrefix(b, "/")
	switch {
	case suffixSlash && prefixSlash:
		return a + b[1:]
	case !suffixSlash && !prefixSlash:
		return a + "/" + b
	}
	return a + b
}
