package rsa

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"hash"
	"io"
	"strings"

	"github.com/zhiyunliu/golibs/bytesconv"
	"github.com/zhiyunliu/golibs/xencoding/base64"
)

type mapItem func() (crypto.Hash, hash.Hash)

var hashFunc map[string]mapItem

func init() {
	hashFunc = map[string]mapItem{}
	hashFunc["sha1"] = func() (crypto.Hash, hash.Hash) {
		return crypto.SHA1, sha1.New()
	}
	hashFunc["sha256"] = func() (crypto.Hash, hash.Hash) {
		return crypto.SHA256, sha256.New()
	}
	hashFunc["md5"] = func() (crypto.Hash, hash.Hash) {
		return crypto.MD5, md5.New()
	}
}

//GenerateKey 生成基于pkcs1的rsa私、公钥对
//bits 密钥位数:1024,2048
func GenerateKey(bits int) (prikey string, pubkey string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	var prikeyBytes = x509.MarshalPKCS1PrivateKey(privateKey)
	var pubkeyBytes = x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	prikey = base64.Encode(prikeyBytes)
	pubkey = base64.Encode(pubkeyBytes)
	return
}

// Encrypt RSA加密
// plainText 加密原串
// publicKey 加密时候用到的公钥
// pkcsType 密钥格式类型:PKCS1,PKCS8
func Encrypt(plainText, publicKey string) (data string, err error) {
	pub, err := getPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("rsa加密时获取公钥err:%v", err)
	}

	res, err := rsa.EncryptPKCS1v15(rand.Reader, pub, bytesconv.StringToBytes(plainText))
	if err != nil {
		return "", fmt.Errorf("ras以[PKCS1v15]加密错误:%v", err)
	}
	data = base64.Encode(res)
	return
}

// Decrypt RSA解密
// ciphertext 解密数据原串
// privateKey 解密时候用到的秘钥
// pkcsType 密钥格式类型:PKCS1,PKCS8
func Decrypt(ciphertext, privateKey string) (data string, err error) {
	priv, err := getPrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("rsa解密时获取私钥失败:%v", err)
	}

	cipherByte, err := base64.Decode(ciphertext)
	if err != nil {
		return "", fmt.Errorf("rsa解密，通过base64解码错误:%v", err)
	}
	res, err := rsa.DecryptPKCS1v15(rand.Reader, priv, cipherByte)
	if err != nil {
		return "", fmt.Errorf("rsa以[PKCS1v15]解密错误:%v", err)
	}

	data = bytesconv.BytesToString(res)
	return
}

// Sign 使用RSA生成签名
// data 签名数据原串
// privateKey 签名时候用到的秘钥
// mode 加密的模式[目前只支持MD5，SHA1，SHA256,不区分大小写]
// pkcsType 密钥格式类型:PKCS1,PKCS8
func Sign(data, privateKey, mode string) (string, error) {
	priv, err := getPrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("rsa Sign getPrivateKey err:%v", err)
	}

	mode = strings.ToLower(mode)
	callback, ok := hashFunc[mode]
	if !ok {
		return "", fmt.Errorf("rsa Sign 无效的Mode:%s,sha1,sha256,md5", mode)
	}
	cryptohash, t := callback()
	io.WriteString(t, data)
	digest := t.Sum(nil)

	signData, err := rsa.SignPKCS1v15(rand.Reader, priv, cryptohash, digest)
	if err != nil {
		return "", err
	}
	return base64.Encode(signData), nil

}

// Verify 校验签名
// src 签名认证数据原串
// sign 签名串
// publicKey 验证签名的公钥
// mode 加密的模式[目前只支持MD5，SHA1，SHA256,不区分大小写]
// pkcsType 密钥格式类型:PKCS1,PKCS8
func Verify(data, sign, publicKey, mode string) (pass bool, err error) {
	//步骤1，加载RSA的公钥
	rsaPub, err := getPublicKey(publicKey)
	if err != nil {
		return false, fmt.Errorf("rsa Verify getPublicKey err:%v", err)
	}

	mode = strings.ToLower(mode)
	callback, ok := hashFunc[mode]

	if !ok {
		return false, fmt.Errorf("rsa Verify 无效的Mode:%s,sha1,sha256,md5", mode)
	}

	signBtytes, err := base64.Decode(sign)

	cryptohash, t := callback()
	io.WriteString(t, data)
	digest := t.Sum(nil)

	err = rsa.VerifyPKCS1v15(rsaPub, cryptohash, digest, signBtytes)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getPrivateKey(privateKey string) (priv *rsa.PrivateKey, err error) {
	prikeyByte, err := base64.Decode(privateKey)
	priv, err = x509.ParsePKCS1PrivateKey(prikeyByte)
	if err != nil {
		err = fmt.Errorf("x509 ParsePKCS1PrivateKey err: %v", err)
	}
	return
}

func getPublicKey(publicKey string) (pub *rsa.PublicKey, err error) {

	pubkeyByte, err := base64.Decode(publicKey)
	pub, err = x509.ParsePKCS1PublicKey(pubkeyByte)
	if err != nil {
		err = fmt.Errorf("x509 ParsePKCS1PublicKey err: %v", err)
	}
	return
}
