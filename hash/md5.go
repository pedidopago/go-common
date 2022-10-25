package hash

import (
	"crypto/md5"
	"fmt"
)

func MD5(s string) []byte {
	h := md5.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func Hex(d []byte) string {
	return fmt.Sprintf("%x", d)
}
