package dao

import (
	"courseSys/types"
	"log"
)

type StateType int

const (
	Admin   types.UserType = 1
	Student types.UserType = 2
	Teacher types.UserType = 3
	Exist   StateType      = 1
	Delete  StateType      = 2
)

type User struct {
	UserID    int64          `gorm:"column:userID;AUTO_INCREMENT;primaryKey"`
	Nickname  string         `gorm:"column:nickName;type:varchar(20);NOT NULL"`
	Username  string         `gorm:"column:userName;type:varchar(20);unique;NOT NULL"`
	Password  string         `gorm:"column:password;type:varchar(20);NOT NULL"`
	UserType  types.UserType `gorm:"column:userType;type:tinyint;NOT NULL"`
	UserState StateType      `gorm:"column:userState;type:tinyint;NOT NULL"`
}

func (User) TableName() string {
	return "user"
}

func userTableInit() {
	// 建user表
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Println("Creating user table failed")
		return
	}
}

func adminRegister() {
	// 提前内置管理员账户
	db.Create(&User{
		Nickname:  "JudgeAdmin",
		Username:  "JudgeAdmin",
		Password:  "JudgePassword2022",
		UserType:  types.Admin,
		UserState: Exist,
	})
}

func UserAuth(username string, password string) (user User) {
	// do not use struct for conditional query
	// select * from user where username=username, password=password
	db.Where("userName = ? AND password = ?", username, password).Find(&user)
	return
}

func UserInsert(user *User) {
	db.Create(user)
}

func UserCheck(username string) (user User) {
	// select * from user where username=username
	db.Where("userName = ?", username).Find(&user)
	return
}

func GetUser(uid int64) (user User) {
	// select * from user where userid=uid
	db.Where("userID = ?", uid).Find(&user)
	return
}

func GetListUser(Limit int, Offset int) (users []User) {
	// select * from user where UserState=1
	db.Where("UserState = ?", 1).Limit(Limit).Offset((Offset - 1) * Limit).Find(&users)
	return
}

func UpdateUserInt(user User, key string, value int) {
	// UPDATE users SET key=value WHERE id=user.UserID;
	db.Model(&user).Update(key, value)
}
func UpdateUserStr(user User, key string, value string) {
	// UPDATE users SET key=value WHERE id=user.UserID;
	db.Model(&user).Update(key, value)
}
