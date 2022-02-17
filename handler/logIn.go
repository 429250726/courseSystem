package handler

import (
	"courseSys/dao"
	"courseSys/types"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
)

func LogInHandler(c *gin.Context){
	// 读取表单参数
	req := &types.LoginRequest{}
	err := c.BindJSON(req)
	if err != nil {
		log.Println("login parameter get fail")
		c.JSON(200, &types.LoginResponse{
			Code: types.ParamInvalid,
			Data: struct {
				UserID string
			}{},
		})
		return
	}
	//username, flag := c.GetPostForm("Username")
	//if flag == false{
	//	log.Println("Login username get fail")
	//}
	//password, flag := c.GetPostForm("Password")
	//if flag == false{
	//	log.Println("Login password get fail")
	//}

	// 用户名和密码检索数据库
	user := dao.UserAuth(req.Username, req.Password)

	// 0代表用户不存在 2代表用户被删除
	if user.UserState != dao.Exist {
		c.JSON(200, &types.LoginResponse{
			Code: types.WrongPassword,
			Data: struct {
    			UserID string
			}{},
		})
		return
	}
	log.Printf("[LogInHandler] user:%+v",user)
	// 设置cookie value为userID,userType 方便查看用户读取
	c.SetCookie(
		"camp-session",
		strconv.FormatInt(user.UserID, 10) + "," + strconv.FormatInt(int64(user.UserType), 10),
		1000,
		"/",
		"180.184.70.138",
		false,
		true,
	)
	//util.DPrintf()
	c.JSON(200, &types.LoginResponse{
		Code: types.OK,
		Data: struct{
			UserID string
		}{
			UserID: strconv.FormatInt(user.UserID, 10),
		},
	})
}

func LogOutHandler(c *gin.Context){
	// 丢弃cookie
	c.SetCookie(
		"camp-session",
		"",
		-1,
		"/",
		"180.184.70.138",
		false,
		true,
	)
	c.JSON(200, &types.LogoutResponse{
		Code: types.OK,
	})
}

func WhoAmIHandler(c *gin.Context){
	// 从登录接口设置的cookie中读取uid
	cookie, err := c.Cookie("camp-session")
	if err != nil{
		log.Println("Get cookie failed")
		c.JSON(
			200,
			&types.WhoAmIResponse{
				Code: types.LoginRequired,
				Data: types.TMember{
				},
			},
		)
		return
	}
	cookieSplit := strings.Split(cookie, ",")
	userID, err := strconv.ParseInt(cookieSplit[0], 10, 64)
	if err != nil{
		log.Println("parseInt fail at UserMsg")
	}
	// 通过uid检索数据库
	user := dao.GetUser(userID)
	c.JSON(
		200,
		&types.WhoAmIResponse{
			Code: types.OK,
			Data: types.TMember{
				UserID:   strconv.FormatInt(user.UserID, 10),
				Nickname: user.Nickname,
				Username: user.Username,
				UserType: user.UserType,
			},
		},
		)

}