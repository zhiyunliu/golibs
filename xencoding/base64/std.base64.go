package base64

import "encoding/base64"

// EncodeBytes 把一个[]byte通过base64编码成string
func Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// DecodeBytes 把一个string通过base64解码成[]byte
func Decode(src string) (s []byte, err error) {
	s, err = base64.StdEncoding.DecodeString(src)
	return
}
