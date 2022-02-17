package dao

import (
	"encoding/json"
	"strconv"
	"time"
)

func WriteCourseToRedis(course Course) {
	rdb.Set("course:"+strconv.FormatInt(course.CourseID, 10), course.CourseCapacity, time.Second)
}
func MinusCap(courseID string) {
	// 原子操作decrease courseCap
	rdb.Decr("course:" + courseID)
}
func IncrCap(courseID string) {
	// 原子操作increase courseCap
	rdb.Incr("course:" + courseID)
}

func RGetAllStudentCourseByStudentID(studentID string) (studentCourses []StudentCourse,err error){
	res,err:=rdb.Get("[RGetAllStudentCourseByStudentID]:"+studentID).Result()
	if err!=nil{
		return nil,err
	}
	err=json.Unmarshal([]byte(res),studentCourses)
	if err!=nil{
		return nil,err
	}
	return studentCourses,nil
}




