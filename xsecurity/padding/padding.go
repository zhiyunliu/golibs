package padding

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	//PaddingNull 不填充
	PaddingNull = "NULL"
	//PaddingPkcs7 .
	PaddingPkcs7 = "PKCS7"
	//PaddingPkcs5 .
	PaddingPkcs5 = "PKCS5"
	//PaddingZero 0填充
	PaddingZero = "ZERO"
)

// ZeroPadding Zero填充模式
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

// ZeroUnPadding 去除Zero的补码
func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

// PKCS5Padding PKCS5填充模式
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding 去除PKCS5的补码
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return origData
	}
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// PKCS7Padding PKCS7填充模式
func PKCS7Padding(data []byte) []byte {
	blockSize := 16
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)

}

// PKCS7UnPadding 去除PKCS7的补码
func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	// 去掉最后一个字节 unpadding 次
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

// GetModePadding 解析加密模式和填充模式
func GetSecretMode(name string) (mode, padding string, err error) {
	names := strings.Split(name, "/")
	if len(names) != 2 {
		err = fmt.Errorf("输入模式不正确:%s", name)
		return
	}
	mode = strings.ToUpper(names[0])
	padding = strings.ToUpper(names[1])

	if strings.EqualFold(padding, PaddingPkcs5) ||
		strings.EqualFold(padding, PaddingPkcs7) ||
		strings.EqualFold(padding, PaddingZero) ||
		strings.EqualFold(padding, PaddingNull) {
		return
	}
	err = fmt.Errorf("填充模式不支持:%s", padding)
	return
}
