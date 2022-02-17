package handler

import (
	"courseSys/dao"
	"courseSys/types"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"strconv"
)

func BookCourse(c *gin.Context){
	req := &types.BookCourseRequest{}
	err := c.BindJSON(req)
	if err != nil {
		log.Println("bookCourse parameter get fail")
		c.JSON(200, &types.BookCourseResponse{
			Code: types.ParamInvalid,
		})
		return
	}

	studentID, err:= strconv.ParseInt(req.StudentID, 10, 64)
	if err != nil{
		log.Println("BookCourse parse studentID failed")
		c.JSON(200, &types.BookCourseResponse{
			Code: types.ParamInvalid,
		})
		return
	}
	courseID, err := strconv.ParseInt(req.CourseID, 10, 64)
	if err != nil {
		log.Println("courseID parse failed in redis")
		c.JSON(200, &types.BookCourseResponse{
			Code: types.ParamInvalid,
		})
		return
	}

	// 检查学生是否已经选该课程
	if check:= dao.CheckStudentCourse(studentID, courseID); check {
		c.JSON(200, &types.BookCourseResponse{Code: types.StudentHasCourse})
		return
	}

	// 取得redis client
	rdb := dao.GetRedis()

	// cache aside处理并发请求
	// 从redis读courseCap
	courseCap, err := rdb.Get("course:" + req.CourseID).Result()

	// redis中course不存在 读db
	if errors.Is(err, redis.Nil) {

		// 读MySQL
		courseDB := dao.GetCourseByCourseID(courseID)
		if courseDB.CourseID==0{
			c.JSON(200, &types.BookCourseResponse{
				Code: types.CourseNotExisted,
			})	//课程不存在
			return
		}

		// 写redis
		dao.WriteCourseToRedis(courseDB)

		// 读redis
		courseCap, err = rdb.Get("course:" + req.CourseID).Result()
		if err!= nil{
			c.JSON(200, &types.BookCourseResponse{
				Code: types.UnknownError,
			})
			return
		}
	}
	courseCapacity, _ := strconv.ParseInt(courseCap, 10, 64)
	if courseCapacity <= 0{
		c.JSON(200, &types.BookCourseResponse{
			Code: types.CourseNotAvailable,
		})
		return
	}

	// 改redis中courseID对应cap
	dao.MinusCap(req.CourseID)

	// 数据库写回  TODO mq 异步写回
	err = dao.UpdateCourseCap(courseID)

	// 写回失败 重试10次
	for i:=0;err != nil && i < 10; i++{
		err = dao.UpdateCourseCap(courseID)
	}
	if err != nil{
		c.JSON(200, &types.BookCourseResponse{
			Code: types.UnknownError,
		})
		dao.IncrCap(req.CourseID)
		return
	}
	dao.StudentCourseInsert(&dao.StudentCourse{
		StudentID: studentID,
		CourseID: courseID,
	})
	// resp
	c.JSON(200, &types.BookCourseResponse{
		Code: types.OK,
	})
}