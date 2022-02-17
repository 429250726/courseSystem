package dao

import (
	"courseSys/util"
	"log"
)

type Bind struct {
	TeacherID int64			`gorm:"column:teacherID;NOT NULL"`
	CourseID int64			`gorm:"column:courseID;NOT NULL"`
}

func (Bind) TableName() string {
	return "bind"
}

func bindTableInit(){
	util.DPrintln("[bindTableInit]")
	// 建bind表
	err := db.AutoMigrate(&Bind{})
	if err != nil {
		log.Fatalf("Creating bind table failed")
		return
	}
}
func GetBindByBind(b *Bind) (bind Bind){
	util.DPrintf("[GetBindByBind], bind:%+v \n",*b)
	db.Where(b).First(&bind)
	return
}

func GetAllBindByTeacherID(tid int64) (binds []Bind){
	util.DPrintln("[GetAllBindByTeacherID], tid:",tid)
	db.Where("teacherID = ?",tid).Find(&binds)
	util.DPrintln("[GetAllBindByTeacherID], res:",binds)
	return
}

func BindCreate(bind *Bind){
	util.DPrintf("[BindCreate] bind:%+v \n",*bind)
	db.Create(bind)
}

func BindDeleteByBind(bind *Bind){
	util.DPrintf("[BindDeleteByBind] bind:%+v \n",*bind)
	db.Where("teacherID = ? AND courseID = ?", bind.TeacherID,bind.CourseID).Delete(Bind{})
	//db.Where(&bind).Delete()
}