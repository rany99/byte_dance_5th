package password

import (
	"crypto/md5"
	"encoding/hex"
)

const Salt string = "WILL1999"

func Md5(str string) string {
	b := []byte(str)
	s := []byte(Salt)
	h := md5.New()
	h.Write(s)
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
