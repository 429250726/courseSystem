package dao

import (
	"courseSys/util"
	"log"
)

type StudentCourse struct {
	StudentID int64 `gorm:"column:studentID;NOT NULL"`
	CourseID  int64 `gorm:"column:courseID;NOT NULL"`
}

func (StudentCourse) TableName() string {
	return "student"
}

func StudentCourseTableInit() {
	// 建StudentCourse表
	err := db.AutoMigrate(&StudentCourse{})
	if err != nil {
		log.Println("Creating StudentCourse table failed")
		return
	}
}
func CheckStudentCourse(sid int64, cid int64) bool{

	check := &StudentCourse{}
	db.Where("studentID = ? AND courseID = ?", sid, cid).First(check)

	// 学生未选该课程
	if check.StudentID == 0 && check.CourseID == 0{
		return false
	}

	// 学生选择了该课程
	return true
}
func GetAllStudentCourseByStudentID(sid int64) (studentCourses []StudentCourse,err error) {
	////先读缓存
	//studentCourses,err=RGetAllStudentCourseByStudentID(strconv.FormatInt(sid,10))
	//if err!=nil{	//缓存未命中
	//	log.Println("[RGetAllStudentCourseByStudentID] cache not hit")
		//读db
		util.DPrintln("[GetAllStudentCourseByStudentID], sid:", sid)
		db.Where("studentID = ?", sid).Find(&studentCourses)
		util.DPrintln("[GetAllBindByTeacherID], res:", studentCourses)
	//	//写回缓存
	//	data,err:=json.Marshal(studentCourses)
	//	if err!=nil{
	//		log.Println("[GetAllStudentCourseByStudentID]json.Marshal(studentCourses) err:",err)
	//		return
	//	}
	//	rdb.Set("RGetAllStudentCourseByStudentID:"+strconv.FormatInt(sid, 10), data, time.Second)
	//	//再从缓存获取
	//	studentCourses,_=RGetAllStudentCourseByStudentID(strconv.FormatInt(sid,10))
	//}
	return
}

func StudentCourseInsert(course *StudentCourse){
	db.Create(course)
}
