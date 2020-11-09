package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var key = []byte("KubeOperator@202")

//@brief: 填充明文
func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

//@brief: 去除填充数据
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//@brief: AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//@brief:AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

//加密
func StringEncrypt(text string) (string, error) {
	pass := []byte(text)
	xpass, err := AesEncrypt(pass, key)
	if err == nil {
		pass64 := base64.StdEncoding.EncodeToString(xpass)
		return pass64, err
	}
	return "", err
}

//解密
func StringDecrypt(text string) (string, error) {
	bytesPass, err := base64.StdEncoding.DecodeString(text)

	if err != nil {
		return "", err
	}

	tpass, err := AesDecrypt(bytesPass, key)
	if err == nil {
		result := string(tpass[:])
		return result, err
	}
	return "", err
}

//func main() {
//	//key的长度必须是16、24或者32字节，分别用于选择AES-128, AES-192, or AES-256
//	var aeskey = []byte("12345678abcdefgh")
//	pass := []byte("vdncloud123456")
//	xpass, err := AesEncrypt(pass, aeskey)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	pass64 := base64.StdEncoding.EncodeToString(xpass)
//	fmt.Printf("加密后:%v\n", pass64)
//
//	bytesPass, err := base64.StdEncoding.DecodeString(pass64)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	tpass, err := AesDecrypt(bytesPass, aeskey)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Printf("解密后:%s\n", tpass)
//}
