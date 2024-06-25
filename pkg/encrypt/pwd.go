package encrypt

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
