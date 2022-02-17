package handler

import (
	"courseSys/dao"
	"courseSys/types"
	"courseSys/util"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreatMember(c *gin.Context) {
	user := dao.User{
		Nickname:  "",
		Username:  "",
		Password:  "",
		UserType:  0,
		UserState: 0,
	}
	req := &types.CreateMemberRequest{}
	err := c.BindJSON(req)
	if err != nil {
		log.Println("creat member get parameter fail")
		c.JSON(200, &types.CreateMemberResponse{
			Code: types.ParamInvalid,
			Data: struct {
    			UserID string
			}{},
		})
		return
	}
	user.Nickname = req.Nickname
	user.Username = req.Username
	user.Password = req.Password
	user.UserType = req.UserType

	// 操作权限检测
	// 从登录接口设置的cookie中读取uid
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		log.Println("Get cookie failed")
		c.JSON(
			200,
			&types.CreateMemberResponse{
				Code: types.LoginRequired,
				Data: struct {
					UserID string
				}{},
			},
		)
		return
	}
	cookieSplit := strings.Split(cookie, ",")
	curUserType, err := strconv.ParseInt(cookieSplit[1], 10, 64)
	if err != nil {
		log.Println("parseInt fail at create member")
	}
	if types.UserType(curUserType) != types.Admin {
		c.JSON(200, &types.CreateMemberResponse{
			Code: types.PermDenied,
			Data: struct {
				UserID string
			}{},
		})
		return
	}

	// 四种参数检查 valid = 1为合法 4次检查位运算和为1即为全部通过
	valid := 1
	valid &= util.CheckNickname(user.Nickname)
	log.Println(valid)
	valid &= util.CheckUserName(user.Username)
	log.Println(valid)
	valid &= util.CheckPassword(user.Password)
	log.Println(valid)
	valid &= util.CheckUserType(user.UserType)
	log.Println(valid)

	// 参数不合法返回
	if valid == 0 {
		c.JSON(200, &types.CreateMemberResponse{
			Code: types.ParamInvalid,
			Data: struct {
				UserID string
			}{},
		})
		return
	}

	// 检查成员是否已经存在 如果存在 state不为0
	tempUser := dao.UserCheck(user.Username)
	if tempUser.UserState != 0 {
		c.JSON(200, &types.CreateMemberResponse{
			Code: types.UserHasExisted,
			Data: struct {
				UserID string
			}{
				UserID: strconv.FormatInt(tempUser.UserID, 10),
			},
		})
		return
	}

	// 插入新用户
	user.UserState = dao.Exist
	dao.UserInsert(&user)
	c.JSON(200, &types.CreateMemberResponse{
		Code: types.OK,
		Data: struct {
			UserID string
		}{
			UserID: strconv.FormatInt(user.UserID, 10),
		},
	})
}

func GetSingleMember(c *gin.Context) {
	/*
		GET请求，参数通过c.Query或c.DefaultQuery获取
	*/
	var req types.GetMemberRequest

	// 获取UserID
	req.UserID = c.DefaultQuery("UserID", "")

	util.DPrintln("[GetSingleMember] UserID:", req.UserID)

	if len(req.UserID) == 0 { //参数不合法，参数为空
		util.DPrintln("[GetSingleMember] len(req.UserID)==0")

		c.JSON(http.StatusOK, &types.GetMemberResponse{
			Code: types.ParamInvalid,
		})
		return
	}

	// 成员不存在state为0，已存在state为1，已删除state为2
	uid, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil { //参数不合法，无法转换为int64
		util.DPrintln("[GetSingleMember] UserID to int64 err:", err)
		c.JSON(http.StatusOK, &types.GetMemberResponse{
			Code: types.ParamInvalid,
		})
		return
	}
	User := dao.GetUser(uid)

	if User.UserState == 1 {
		c.JSON(200, &types.GetMemberResponse{
			Code: types.OK,
			Data: types.TMember{
				UserID:   strconv.FormatInt(User.UserID, 10),
				Nickname: User.Nickname,
				Username: User.Username,
				UserType: User.UserType,
			},
		})
		return
	}
	if User.UserState == 0 {
		c.JSON(200, &types.GetMemberResponse{
			Code: types.UserNotExisted,
			Data: types.TMember{},
		})
		return
	}
	if User.UserState == 2 {
		c.JSON(200, &types.GetMemberResponse{
			Code: types.UserHasDeleted,
			Data: types.TMember{},
		})
		return
	}
}

func GetMemberList(c *gin.Context) {
	var getMemberListRequest types.GetMemberListRequest

	// 获取Offset与Limit数值
	var tempStr string
	tempStr = c.DefaultQuery("Offset", "")
	util.DPrintln("[GetMemberList] Offset:", tempStr)
	if len(tempStr) == 0 { //参数不合法，参数为空
		util.DPrintln("[GetMemberList] len(Offset)==0")

		c.JSON(200, &types.GetMemberListResponse{
			Code: types.ParamInvalid,
			Data: struct {
				MemberList []types.TMember
			}{},
		})
		return
	}
	getMemberListRequest.Offset, _ = strconv.Atoi(tempStr)

	tempStr = c.DefaultQuery("Limit", "")
	util.DPrintln("[GetMemberList] Limit:", tempStr)
	if len(tempStr) == 0 { //参数不合法，参数为空
		util.DPrintln("[GetMemberList] len(Limit)==0")

		c.JSON(200, &types.GetMemberListResponse{
			Code: types.ParamInvalid,
			Data: struct {
				MemberList []types.TMember
			}{},
		})
		return
	}
	getMemberListRequest.Limit, _ = strconv.Atoi(tempStr)

	// 获取User List
	Users := dao.GetListUser(getMemberListRequest.Limit, getMemberListRequest.Offset)
	// log.Println(Users)

	var MemberList []types.TMember
	for _, User := range Users {
		var tempTMember types.TMember
		tempTMember.UserID = strconv.FormatInt(User.UserID, 10)
		tempTMember.Username = User.Username
		tempTMember.Nickname = User.Nickname
		tempTMember.UserType = User.UserType
		MemberList = append(MemberList, tempTMember)
	}

	c.JSON(200, &types.GetMemberListResponse{
		Code: types.OK,
		Data: struct {
			MemberList []types.TMember
		}{MemberList},
	})
}

func DeleteMember(c *gin.Context) {
	// 获取UserID
	req := &types.DeleteMemberRequest{}
	err := c.BindJSON(req)
	if err != nil {
		log.Println("Delete member parameter get fail")
		c.JSON(200, &types.DeleteMemberResponse{Code: types.ParamInvalid})
		return
	}

	// 成员不存在state为0，已存在state为1，已删除state为2
	tempID, _ := strconv.ParseInt(req.UserID, 10, 64)
	User := dao.GetUser(tempID)

	if User.UserState == 1 {
		// 更新user的state为2
		dao.UpdateUserInt(User, "UserState", 2)

		c.JSON(200, &types.DeleteMemberResponse{
			Code: types.OK,
		})
		return
	}
	if User.UserState == 0 {
		c.JSON(200, &types.GetMemberResponse{
			Code: types.UserNotExisted,
		})
		return
	}
	if User.UserState == 2 {
		c.JSON(200, &types.GetMemberResponse{
			Code: types.UserHasDeleted,
		})
		return
	}
}

func UpdateMember(c *gin.Context) {
	req := &types.UpdateMemberRequest{}
	err := c.BindJSON(req)
	if err != nil {
		log.Println("update member get parameter fail")
		c.JSON(200, &types.UpdateMemberResponse{Code: types.ParamInvalid})
		return
	}

	// 检查Nickname是否合法
	valid := 1
	valid &= util.CheckNickname(req.Nickname)
	log.Println(valid)
	// 参数不合法返回
	if valid == 0 {
		c.JSON(200, &types.CreateMemberResponse{
			Code: types.ParamInvalid,
			Data: struct {
				UserID string
			}{},
		})
		return
	}

	// 成员不存在state为0，已存在state为1，已删除state为2
	tempID, _ := strconv.ParseInt(req.UserID, 10, 64)
	User := dao.GetUser(tempID)

	if User.UserState == 1 {
		// 更新user的昵称
		dao.UpdateUserStr(User, "Nickname", req.Nickname)

		c.JSON(200, &types.UpdateMemberResponse{
			Code: types.OK,
		})
		return
	}
	if User.UserState == 0 {
		c.JSON(200, &types.UpdateMemberResponse{
			Code: types.UserNotExisted,
		})
		return
	}
	if User.UserState == 2 {
		c.JSON(200, &types.UpdateMemberResponse{
			Code: types.UserHasDeleted,
		})
		return
	}
}
