package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
)

func IntMix(a, b int) int {
	if a+b < 10 {
		return a + b
	}
	if a < b {
		return a
	}
	return b
}

func main1() {
	fmt.Print(IntMix(1, 2))
	fmt.Print(IntMix(1, 2))
}

type MockR struct {
	c int8
}

func (b *MockR) Read(p []byte) (n int, err error) {
	rand.Read(p)

	return len(p), nil
}

// pkcs7strip remove pkcs7 padding
func pkcs7strip(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs7: Data is empty")
	}
	if length%blockSize != 0 {
		return nil, errors.New("pkcs7: Data is not block-aligned")
	}
	padLen := int(data[length-1])
	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
	if padLen > blockSize || padLen == 0 || !bytes.HasSuffix(data, ref) {
		return nil, errors.New("pkcs7: Invalid padding")
	}
	return data[:length-padLen], nil
}

func pkcs7pad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 1 || blockSize >= 256 {
		return nil, fmt.Errorf("pkcs7: Invalid block size %d", blockSize)
	} else {
		padLen := blockSize - len(data)%blockSize
		padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
		return append(data, padding...), nil
	}
}

func main() {
	var mr MockR
	message := "===abc==AAA"
	fmt.Println("message: ", string(message))

	//rand.Reader
	//生成私钥
	privateKey, err := rsa.GenerateKey(&mr, 1024)
	if err != nil {
		panic(err)
	}

	//生成公钥
	publicKey := privateKey.PublicKey
	//fmt.Println(publicKey)

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA1, []byte(message), nil)
	if err != nil {
		fmt.Println("SignPSS fail:", err)
		return
	}
	fmt.Printf("%x\n", signature)

	//根据公钥加密
	encryptedBytes, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, &publicKey, []byte(message), nil)
	if err != nil {
		fmt.Println("EncryptOAEP fail:", err)
		return
	}
	//fmt.Println("encryptedBytes: ", encryptedBytes)

	// 加密后进行base64编码
	encryptBase64 := base64.StdEncoding.EncodeToString(encryptedBytes)
	fmt.Printf("encrypted(%d): %v\n", len(encryptedBytes), encryptBase64)

	// 解密
	// base64解码
	decodedBase64, err := base64.StdEncoding.DecodeString(encryptBase64)
	if err != nil {
		fmt.Println("decode fail:", err)
		return
	}
	//fmt.Println("decodedBase64: ", decodedBase64)

	//根据私钥解密
	decryptedBytes, err := privateKey.Decrypt(nil, decodedBase64, &rsa.OAEPOptions{Hash: crypto.SHA1})
	if err != nil {
		fmt.Println("decrypt fail:", err)
		return
	}
	fmt.Println("decrypted:", string(decryptedBytes))

}
