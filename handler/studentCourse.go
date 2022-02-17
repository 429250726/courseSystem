package handler

import (
	"courseSys/dao"
	"courseSys/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetStudentCourseHandler(c *gin.Context) {
	req := &types.GetStudentCourseRequest{}
	resp := &types.GetStudentCourseResponse{}
	//set param
	req.StudentID = c.DefaultQuery("StudentID", "")
	if len(req.StudentID) == 0 {
		resp.Code = types.StudentNotExisted
		c.JSON(http.StatusOK, resp)
		return
	}
	//
	GetStudentCourse(req, resp)
	c.JSON(http.StatusOK, resp)
}
func GetStudentCourse(req *types.GetStudentCourseRequest, resp *types.GetStudentCourseResponse) {
	/*
		获取StudentID拥有的所有课程
	*/
	log.Println("[GetStudentCourse] req:", *req)

	StudentID, err := strconv.ParseInt(req.StudentID, 10, 64)
	if err != nil {
		log.Printf("[GetStudentCourse] StudentID to int64 err,StudentID:%v, err:%v \n", req.StudentID, err)
		resp.Code = types.ParamInvalid
		return
	}
	studentCourses ,err:= dao.GetAllStudentCourseByStudentID(StudentID)
	if err!=nil{
		resp.Code=types.UnknownError
		return
	}
	if len(studentCourses) == 0 { //学生没有课程
		resp.Code = types.StudentHasNoCourse
		return
	}
	resp.Code = types.OK
	cnt := len(studentCourses)
	resp.Data.CourseList = make([]types.TCourse, cnt)
	for i := 0; i < cnt; i++ {
		course, bind := dao.GetCourse(studentCourses[i].CourseID)
		resp.Data.CourseList[i] = types.TCourse{
			CourseID:  strconv.FormatInt(studentCourses[i].CourseID, 10),
			Name:      course.CourseName,
			TeacherID: strconv.FormatInt(bind.TeacherID, 10),
		}
	}
}
