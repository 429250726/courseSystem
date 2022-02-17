package dao

import (
	"courseSys/util"
	"gorm.io/gorm/clause"
	"log"
	"strconv"
)

type Course struct {
	CourseID int64 			`gorm:"column:courseID;AUTO_INCREMENT;primaryKey"`
	CourseName string		`gorm:"column:courseName;varchar(20);NOT NULL"`
	CourseCapacity int64	`gorm:"column:courseCapacity;int;NOT NULL"`
}
func (Course) TableName() string {
	return "course"
}

func courseTableInit(){
	// 建course表
	err := db.AutoMigrate(&Course{})
	if err != nil {
		log.Println("Creating course table failed")
		return
	}
}

func GetCourseByCourseID(cid int64) (course Course){
	db.Where("courseID = ?", cid).First(&course)
	return
}

func CourseCreate(course *Course){
	db.Create(course)
}

func GetCourse(cid int64)(course Course, bind Bind){
	util.DPrintln("[GetCourse] cid:",cid)
	db.Where("courseID = ?", cid).First(&course)
	db.Where("courseID = ?", cid).First(&bind)
	return
}

func UpdateCourseCap(cid int64) error{
	// db事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	course := Course{}
	err := tx.First(&course, Course{CourseID: cid}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 锁
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&course)

	course.CourseCapacity--
	if err := tx.Model(&course).Update("courseCapacity", course.CourseCapacity).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 更新DB事务后 再删除缓存
	rdb.Del("courseID:" + strconv.FormatInt(cid, 10))
	return nil
}