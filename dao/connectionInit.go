package dao

import (
	"log"

	"github.com/go-redis/redis"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var rdb *redis.Client

func DBInit() {
	var err error
	//db, err = gorm.Open(mysql.Open("root:429250726@tcp(127.0.0.1:3306)/camp?charset=utf8"), &gorm.Config{})
	db, err = gorm.Open(mysql.Open("root:bytedancecamp@tcp(180.184.70.138)/camp?charset=utf8"), &gorm.Config{})
	if err != nil {
		log.Fatalf("mysql connection failed")
		return
	}

	// 初始化用户表
	userTableInit()

	// 提前内置管理员账户
	adminRegister()

	//初始化course表
	courseTableInit()

	//初始化bind表
	bindTableInit()

	//初始化CourseStudent表
	StudentCourseTableInit()
}
func RedisInit() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "180.184.70.138:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Println("redis init fail")
	} else {
		log.Println("redis init succeed")
	}
}
func GetRedis() *redis.Client {
	return rdb
}
