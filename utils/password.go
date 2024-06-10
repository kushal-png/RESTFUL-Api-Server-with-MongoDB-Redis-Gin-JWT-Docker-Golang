package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string)(string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func VerifyPassword(pw string, hashPw string)(error){
	err:=bcrypt.CompareHashAndPassword([]byte(hashPw),[]byte(pw))
	if err!=nil{
		return err
	}
	return nil
}
