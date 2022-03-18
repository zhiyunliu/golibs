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
	plainBytes := dataPadding(plainText, p, opts.BlockSize)

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

	cipherBytes := dataUnpadding(cipherText, p)

	plainBytes, err := decryptData(aesCipher, cipherBytes, opts.IV, encMode)
	if err != nil {
		return "", err
	}
	plainText = strings.TrimSpace(string(plainBytes))
	return
}

func dataPadding(plainText string, padingmode string, blockSize int) (plainBytes []byte) {
	switch padingmode {
	case padding.PaddingPkcs7:
		plainBytes = padding.PKCS7Padding(bytesconv.StringToBytes(plainText))
	case padding.PaddingPkcs5:
		plainBytes = padding.PKCS5Padding(bytesconv.StringToBytes(plainText), blockSize)
	case padding.PaddingZero:
		plainBytes = padding.ZeroPadding(bytesconv.StringToBytes(plainText), blockSize)
	default:
		plainBytes = bytesconv.StringToBytes(plainText)
	}
	return
}

func dataUnpadding(cipherText string, padingmode string) (cipherBytes []byte) {
	switch padingmode {
	case padding.PaddingPkcs7:
		cipherBytes = padding.PKCS7UnPadding(bytesconv.StringToBytes(cipherText))
	case padding.PaddingPkcs5:
		cipherBytes = padding.PKCS5UnPadding(bytesconv.StringToBytes(cipherText))
	case padding.PaddingZero:
		cipherBytes = padding.ZeroUnPadding(bytesconv.StringToBytes(cipherText))
	default:
		cipherBytes = bytesconv.StringToBytes(cipherText)
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
