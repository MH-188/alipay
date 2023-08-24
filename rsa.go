package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

// Rsa2PrivateSign RSA2私钥签名
func Rsa2PrivateSign(signContent string, privateKey *rsa.PrivateKey, hash crypto.Hash) string {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashedContent := shaNew.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, hash, hashedContent)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}

// RSAVerifyWithKey 验证签名
func RSAVerifyWithKey(content []byte, sign []byte, publicKey *rsa.PublicKey, hash crypto.Hash) error {
	signBytes, err := base64.StdEncoding.DecodeString(string(sign))
	if err != nil {
		return err
	}
	var h = hash.New()
	h.Write(content)
	var hashed = h.Sum(nil)

	return rsa.VerifyPKCS1v15(publicKey, hash, hashed, signBytes)
}

// ParsePrivateKey 解析私钥
func ParsePrivateKey(environment string, privateKeyBytes []byte) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("私钥信息错误！")
	}

	var key interface{}
	var err error
	if environment == CLIENT_ENVIRONMENT_PROD {
		//生产环境
		key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	} else {
		//沙箱环境生成的公私钥
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	pKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%s 不是 rsa private key", string(privateKeyBytes))
	}
	return pKey, nil
}

// ParseCertificate 解析证书
func ParseCertificate(b []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("证书信息错误！")
	}
	csr, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return csr, nil
}

// ParsePublicKey 解析公钥
func ParsePublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("传入文件不是 rsa public key")
	}

	return key, err
}

// LoadPrivateKeyWithPath 通过私钥的文件路径内容加载私钥
func LoadPrivateKeyWithPath(environment string, path string) (privateKey *rsa.PrivateKey, err error) {
	privateKeyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件错误: %s", err.Error())
	}
	return ParsePrivateKey(environment, privateKeyBytes)
}

// LoadPublicKeyWithPath 通过公钥的文件路径内容加载公钥
func LoadPublicKeyWithPath(path string) (privateKey *rsa.PublicKey, err error) {
	privateKeyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件错误: %s", err.Error())
	}
	return ParsePublicKey(privateKeyBytes)
}
