package main

import (
	"courseSys/dao"
	"courseSys/handler"
	"courseSys/util"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func main() {
	// 初始化数据库连接，建表
	dao.DBInit()

	// 初始化redis连接
	dao.RedisInit()

	r := gin.Default()
	g := r.Group("/api/v1")

	// ping测试
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong!",
		})
	})
	// 成员管理
	g.POST("/member/create", handler.CreatMember)
	g.GET("/member", handler.GetSingleMember)
	g.GET("/member/list", handler.GetMemberList)
	g.POST("/member/update", handler.UpdateMember)
	g.POST("/member/delete", handler.DeleteMember)

	// 登录
	g.POST("/auth/login", handler.LogInHandler)
	g.POST("/auth/logout", handler.LogOutHandler)
	g.GET("/auth/whoami", handler.WhoAmIHandler)

	// 排课
	g.POST("/course/create", handler.CourseCreate)
	g.GET("/course/get", handler.GetCourse)

	g.POST("/teacher/bind_course", handler.BindCourseHandler)
	g.POST("/teacher/unbind_course", handler.UnbindCourseHandler)
	g.GET("/teacher/get_course", handler.GetTeacherCourseHandler)
	g.POST("/course/schedule", handler.ScheduleCourseHandler)

	// 抢课
	// 限流桶 500 token/per second  500 token maximum
	limit := util.NewLimiter(500, 500, time.Second)
	g.POST("/student/book_course", limit, handler.BookCourse)
	g.GET("/student/course", handler.GetStudentCourseHandler)

	if len(os.Args) > 1 { //方便启动多个服务
		r.Run(":" + os.Args[1])
	} else {
		r.Run(":80") //camp要求在80端口启动服务
	}
}
