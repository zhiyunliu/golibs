package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"strings"

	"github.com/zhiyunliu/golibs/bytesconv"
	"github.com/zhiyunliu/golibs/xencoding/base64"
	"github.com/zhiyunliu/golibs/xsecurity/padding"
)

// mode: cbc/pkcs7
func Encrypt(plainText string, key string, mode string, opt ...Option) (cipherText string, err error) {
	keyBytes := []byte(key)

	aesCipher, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("aes.Encrypt.NewCipher:%v", err)
	}

	opts := options{
		IV:        make([]byte, 16),
		BlockSize: aes.BlockSize,
	}
	for i := range opt {
		opt[i](&opts)
	}

	encMode, p, err := padding.GetSecretMode(mode)
	if err != nil {
		return "", err
	}
	plainBytes := dataPadding(bytesconv.StringToBytes(plainText), p, opts.BlockSize)

	cipherBytes, err := encryptData(aesCipher, plainBytes, opts.IV, encMode)
	if err != nil {
		return "", err
	}

	cipherText = base64.Encode(cipherBytes)
	return
}

func Decrypt(cipherText string, key string, mode string, opt ...Option) (plainText string, err error) {

	keyBytes := []byte(key)
	aesCipher, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("aes.Decrypt.NewCipher:%v", err)
	}

	opts := options{
		IV:        make([]byte, 16),
		BlockSize: aes.BlockSize,
	}
	for i := range opt {
		opt[i](&opts)
	}

	encMode, p, err := padding.GetSecretMode(mode)
	if err != nil {
		return "", err
	}

	cipherBytes, err := base64.Decode(cipherText)
	if err != nil {
		return
	}

	plainBytes, err := decryptData(aesCipher, cipherBytes, opts.IV, encMode)
	if err != nil {
		return "", err
	}
	plainBytes, err = dataUnpadding(plainBytes, p)
	if err != nil {
		return "", err
	}
	plainText = strings.TrimSpace(string(plainBytes))
	return
}

func dataPadding(plainBytes []byte, padingmode string, blockSize int) (paddingBytes []byte) {
	switch padingmode {
	case padding.PaddingPkcs7:
		paddingBytes = padding.PKCS7Padding(plainBytes)
	case padding.PaddingPkcs5:
		paddingBytes = padding.PKCS5Padding(plainBytes, blockSize)
	case padding.PaddingZero:
		paddingBytes = padding.ZeroPadding(plainBytes, blockSize)
	default:
		paddingBytes = plainBytes
	}
	return
}

func dataUnpadding(plainBytes []byte, padingmode string) (cipherBytes []byte, err error) {

	switch padingmode {
	case padding.PaddingPkcs7:
		cipherBytes = padding.PKCS7UnPadding(plainBytes)
	case padding.PaddingPkcs5:
		cipherBytes = padding.PKCS5UnPadding(plainBytes)
	case padding.PaddingZero:
		cipherBytes = padding.ZeroUnPadding(plainBytes)
	default:
		cipherBytes = plainBytes
	}
	return
}

func encryptData(aespher cipher.Block, plainBytes []byte, iv []byte, encMode string) (cipherBytes []byte, err error) {

	cipherBytes = make([]byte, len(plainBytes))
	var stream cipher.Stream

	switch {
	case strings.EqualFold(encMode, "cbc"):
		blockMode := cipher.NewCBCEncrypter(aespher, iv)
		blockMode.CryptBlocks(cipherBytes, plainBytes)
	case strings.EqualFold(encMode, "ctr"):
		stream = cipher.NewCTR(aespher, iv)
		stream.XORKeyStream(cipherBytes, plainBytes)
	case strings.EqualFold(encMode, "cfb"):
		stream = cipher.NewCFBEncrypter(aespher, iv)
		stream.XORKeyStream(cipherBytes, plainBytes)
	case strings.EqualFold(encMode, "ofb"):
		stream = cipher.NewOFB(aespher, iv)
		stream.XORKeyStream(cipherBytes, plainBytes)
	default:
		return nil, fmt.Errorf("aes.encrypt.不支持的加密模式:%v", encMode)
	}
	return
}

func decryptData(aespher cipher.Block, cipherBytes []byte, iv []byte, encMode string) (plainBytes []byte, err error) {

	plainBytes = make([]byte, len(cipherBytes))
	var stream cipher.Stream

	switch {
	case strings.EqualFold(encMode, "cbc"):
		blockMode := cipher.NewCBCDecrypter(aespher, iv)
		blockMode.CryptBlocks(plainBytes, cipherBytes)
	case strings.EqualFold(encMode, "ctr"):
		stream = cipher.NewCTR(aespher, iv)
		stream.XORKeyStream(plainBytes, cipherBytes)
	case strings.EqualFold(encMode, "cfb"):
		stream = cipher.NewCFBDecrypter(aespher, iv)
		stream.XORKeyStream(plainBytes, cipherBytes)
	case strings.EqualFold(encMode, "ofb"):
		stream = cipher.NewOFB(aespher, iv)
		stream.XORKeyStream(plainBytes, cipherBytes)
	default:
		return nil, fmt.Errorf("aes.decrypt.不支持的加密模式:%v", encMode)
	}
	return
}
