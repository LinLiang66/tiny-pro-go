package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

// Encry 使用 PBKDF2 算法进行加密
func Encry(value, salt string) (string, error) {
	// 定义 PBKDF2 的参数
	iterations := 1000 // 迭代次数
	keyLength := 18    // 输出密钥的长度（以字节为单位）

	// 使用 PBKDF2 生成密钥
	hashedPassword := pbkdf2.Key([]byte(value), []byte(salt), iterations, keyLength, sha256.New)

	// 将字节数组转换为 Base64 编码的字符串
	return base64.StdEncoding.EncodeToString(hashedPassword), nil
}

// GenerateSalt 生成随机盐值
func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// VerifyPassword 验证密码是否匹配
func VerifyPassword(password, salt, hash string) (bool, error) {
	computedHash, err := Encry(password, salt)
	if err != nil {
		return false, err
	}
	// 使用 subtle.ConstantTimeCompare 进行安全比较
	return subtle.ConstantTimeCompare([]byte(computedHash), []byte(hash)) == 1, nil
}
