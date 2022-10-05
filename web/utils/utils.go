package utils

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	salt = "pwd_salt"
)

// Pwd2Hash 密码转换成hash
func Pwd2Hash(pwd string) string {
	m5 := md5.New()
	m5.Write([]byte(pwd))
	pwd_hash := hex.EncodeToString(m5.Sum([]byte(salt)))
	return pwd_hash
}
