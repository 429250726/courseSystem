package util

import (
	"courseSys/types"
	"unicode"
)

func CheckUserName(str string) (valid int){
	if len(str) < 8 || len(str) > 20 {
		return
	}
	for _, char := range str{
		if !unicode.IsNumber(char) && !unicode.IsLetter(char){
			return
		}
	}
	valid = 1 // 合法
	return
}
func CheckNickname(nickname string) (valid int){
	if len(nickname) >= 4 && len(nickname) <= 20 {
		valid = 1 // 合法
	}
	return
}
func CheckPassword(str string) (valid int){
	if len(str) < 8 || len(str) > 20 {
		return
	}
	uppercase := 0
	lowercase := 0
	number := 0
	other := 0
	for i := 0; i < len(str); i++ {
		switch {
		case 64 < str[i] && str[i] < 91:
			uppercase += 1
		case 96 < str[i] && str[i] < 123:
			lowercase += 1
		case 47 < str[i] && str[i] < 58:
			number += 1
		default:
			other += 1
		}
	}
	if other != 0{
		return
	}
	if uppercase != 0 && lowercase != 0 && number != 0 {
		valid = 1 // 合法
	}
	return
}
func CheckUserType(userType types.UserType) (valid int){
	if userType == types.Admin || userType == types.Student || userType == types.Teacher{
		valid = 1 // 合法
	}
	return
}
