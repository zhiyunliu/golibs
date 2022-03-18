package hex

import (
	"encoding/hex"
)

// Encode 把[]byte类型通过hex编码成string
func Encode(src []byte) string {
	return hex.EncodeToString(src)
}

// Decode 把一个string类型通过hex解码成string
func Decode(src string) (r []byte, err error) {
	return hex.DecodeString(src)
}
