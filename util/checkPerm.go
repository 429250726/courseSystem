package util

import (
	"courseSys/types"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
)

//CheckPerm 判断用户是否是指定类型，用于鉴权
func CheckPerm(c *gin.Context,targetType []types.UserType) (is bool,errNo types.ErrNo){
	// 从登录接口设置的cookie中读取uid
	cookie, err := c.Cookie("camp-session")
	if err != nil{
		log.Println("Get cookie failed")
		return false, types.LoginRequired	//未登录
	}
	cookieSplit := strings.Split(cookie, ",")
	curUserType, err := strconv.ParseInt(cookieSplit[1], 10, 64)
	if err != nil{
		log.Println("userType must be integer")
		return false,types.ParamInvalid	//参数无效
	}
	for _,tp:=range targetType{
		if types.UserType(curUserType)==tp{
			return true,types.OK
		}
	}
	return false,types.PermDenied	//权限不够
}

