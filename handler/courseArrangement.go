package handler

import (
	"courseSys/dao"
	"courseSys/types"
	"courseSys/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func CourseCreate(c *gin.Context){
	// 读取request参数
	req := &types.CreateCourseRequest{}
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("courseCreat parameter get fail")
		c.JSON(200, &types.CreateCourseResponse{
			Code: types.ParamInvalid,
			Data: struct {
				CourseID string
			}{},
		})
		return
	}

	// new course
	course := &dao.Course{
		CourseName:     req.Name,
		CourseCapacity: int64(req.Cap),
	}

	// 插入数据库
	dao.CourseCreate(course)

	// 获取返回的courseID
	courseID := strconv.Itoa(int(course.CourseID))
	c.JSON(200, &types.CreateCourseResponse{
		Code: types.OK,
		Data: struct {
    		CourseID string
		}{
			CourseID: courseID,
		},
	})
}

func GetCourse(c *gin.Context){

	// 读取参数
	cid := c.DefaultQuery("CourseID", "")
	courseID, err:= strconv.ParseInt(cid, 10, 64)
	if err != nil {
		log.Println("GetCourse parse cid fail")
	}

	// 查询课程信息 绑定信息
	course, bind := dao.GetCourse(courseID)
	c.JSON(200, &types.GetCourseResponse{
		Code: types.OK,
		Data: types.TCourse{
			CourseID:  strconv.FormatInt(course.CourseID, 10),
			Name:      course.CourseName,
			TeacherID: strconv.FormatInt(bind.TeacherID, 10),
		},
	})
}

//

func BindCourseHandler(c *gin.Context){
	req:=&types.BindCourseRequest{}
	resp:=&types.BindCourseResponse{}
	//check perm
	ok,errNo:=util.CheckPerm(c,[]types.UserType{types.Teacher,types.Admin})
	if !ok{
		resp.Code=errNo
		log.Printf("[BindCourseHandler] resp:%+v",*resp)
		c.JSON(200,resp)
		return
	}
	//get param
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("Bind course parameter get fail")
		resp.Code=types.ParamInvalid
		log.Printf("[BindCourseHandler] resp:%+v",*resp)
		c.JSON(http.StatusOK,resp)
		return
	}


	BindCourse(req,resp)
	log.Printf("[BindCourseHandler] resp:%+v",*resp)
	c.JSON(http.StatusOK,resp)
}
func BindCourse(req *types.BindCourseRequest,resp *types.BindCourseResponse){
	/*
	TeacherID不需要做已落库校验
	CourseID需要判断是否存在、是否已绑定
	*/
	log.Printf("[BindCourse] req:%+v",*req)

	CourseID,err:=strconv.ParseInt(req.CourseID,10,64)
	if err!=nil{
		log.Printf("[BindCourse] CourseID to int64 err, CourseID:%v, err:%v \n",req.CourseID,err)
		resp.Code=types.ParamInvalid
		return
	}
	TeacherID,err:=strconv.ParseInt(req.TeacherID,10,64)
	if err!=nil{
		log.Printf("[BindCourse] TeacherID to int64 err, TeacherID:%v, err:%v \n",req.TeacherID,err)
		resp.Code=types.ParamInvalid
		return
	}
	//判断课程是否存在
	Course:=dao.GetCourseByCourseID(CourseID)
	if Course.CourseID==0{
		log.Printf("[BindCourse] CourseNotExisted,CourseID:%v \n",req.CourseID)
		resp.Code=types.CourseNotExisted	//课程不存在
		return
	}
	//判断课程是否已被绑定
	bind:=dao.GetBindByBind(&dao.Bind{CourseID: CourseID})
	if bind.TeacherID!=0{
		log.Printf("[BindCourse] CourseHasBound,CourseID:%v \n",req.CourseID)
		resp.Code=types.CourseHasBound	//课程已被绑定
		return
	}
	bind.TeacherID=TeacherID
	bind.CourseID=CourseID
	dao.BindCreate(&bind)
}

func UnbindCourseHandler(c *gin.Context){
	req:=&types.UnbindCourseRequest{}
	resp:=&types.UnbindCourseResponse{}
	//check perm
	ok,errNo:=util.CheckPerm(c,[]types.UserType{types.Teacher,types.Admin})
	if !ok{
		resp.Code=errNo
		log.Printf("[UnbindCourseHandler] resp:%+v",*resp)
		c.JSON(200,resp)
		return
	}
	//
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("UnBind course parameter get fail")
		resp.Code=types.ParamInvalid
		log.Printf("[UnbindCourseHandler] resp:%+v",*resp)
		c.JSON(http.StatusOK,resp)
		return
	}

	//
	UnbindCourse(req,resp)
	log.Printf("[UnbindCourseHandler] resp:%+v",*resp)
	c.JSON(http.StatusOK,resp)
}
func UnbindCourse(req *types.UnbindCourseRequest,resp *types.UnbindCourseResponse){
	/*
		必须已存在TeacherID与CourseID的bind
	*/
	log.Printf("[UnbindCourse] req:%+v",*req)

	CourseID,err:=strconv.ParseInt(req.CourseID,10,64)
	if err!=nil{
		log.Printf("ParamInvalid,CourseID:%v \n",req.CourseID)
		resp.Code=types.ParamInvalid
		return
	}
	TeacherID,err:=strconv.ParseInt(req.TeacherID,10,64)
	if err!=nil{
		log.Printf("ParamInvalid,TeacherID:%v \n",req.TeacherID)
		resp.Code=types.ParamInvalid
		return
	}
	bind:=dao.GetBindByBind(&dao.Bind{CourseID: CourseID})
	util.DPrintf("[UnbindCourse] res of GetBindByBind:%+v",bind)

	if bind.CourseID==0{
		log.Println("[UnbindCourse] CourseID NotExist,CourseID:",CourseID)
		resp.Code=types.CourseNotBind	//CourseID与TeacherID未绑定
		return
	}
	if bind.TeacherID!=TeacherID{
		log.Printf("[UnbindCourse] CourseID is not bind to TeacherID,CourseID:%v,TeacherID:%v \n",req.CourseID,req.TeacherID)
		resp.Code=types.ParamInvalid	//CourseID与TeacherID未绑定，参数不合法
		return
	}
	dao.BindDeleteByBind(&bind)
}
func GetTeacherCourseHandler(c *gin.Context){
	/*
		http.GET
	*/
	req:=&types.GetTeacherCourseRequest{}
	resp:=&types.GetTeacherCourseResponse{}
	//check perm
	ok,errNo:=util.CheckPerm(c,[]types.UserType{types.Teacher,types.Admin})
	if !ok{
		resp.Code=errNo
		log.Printf("[GetTeacherCourseHandler] resp:%+v",*resp)
		c.JSON(200,resp)
		return
	}
	//set param
	req.TeacherID=c.DefaultQuery("TeacherID","")
	if len(req.TeacherID)==0{
		resp.Code=types.ParamInvalid
		log.Printf("[GetTeacherCourseHandler] resp:%+v",*resp)
		c.JSON(http.StatusOK,resp)
		return
	}
	//
	GetTeacherCourse(req,resp)
	log.Printf("[GetTeacherCourseHandler] resp:%+v",*resp)
	c.JSON(http.StatusOK,resp)
}
func GetTeacherCourse(req *types.GetTeacherCourseRequest,resp *types.GetTeacherCourseResponse){
	/*
		获取TeacherID绑定的所有课程
	*/
	log.Printf("[GetTeacherCourse] req:%+v",*req)

	TeacherID,err:=strconv.ParseInt(req.TeacherID,10,64)
	if err!=nil{
		log.Printf("[GetTeacherCourse] TeacherID to int64 err,TeacherID:%v,err:%v \n",req.TeacherID,err)
		resp.Code=types.ParamInvalid
		return
	}
	binds:=dao.GetAllBindByTeacherID(TeacherID)
	resp.Code=types.OK
	cnt:=len(binds)
	resp.Data.CourseList=make([]*types.TCourse,cnt)
	for i:=0;i<cnt;i++{
		name:=dao.GetCourseByCourseID(binds[i].CourseID).CourseName
		resp.Data.CourseList[i]=&types.TCourse{
			CourseID: strconv.FormatInt(binds[i].CourseID,10),
			Name: name,
			TeacherID: strconv.FormatInt(binds[i].TeacherID,10),
		}
	}
}

func ScheduleCourseHandler(c *gin.Context){
	req:=&types.ScheduleCourseRequest{}
	resp:=&types.ScheduleCourseResponse{}
	err := c.BindJSON(req)
	if err != nil {
		log.Printf("Schedule course get parameter fail")
		resp.Code = types.ParamInvalid
		c.JSON(http.StatusOK, resp)
		return
	}

	ScheduleCourse(req,resp)
	log.Printf("[ScheduleCourseHandler] resp:%+v",*resp)
	c.JSON(http.StatusOK,resp)
}
func ScheduleCourse(req *types.ScheduleCourseRequest,resp *types.ScheduleCourseResponse){
	/*
		排课求解器，使老师绑定课程的最优解， 老师有且只能绑定一个课程
		二分图匹配，匹配结果存进resp中
	*/
	log.Printf("[ScheduleCourse] req:%+v",*req)

	//TODO teacherID和courseID的检测
	//TODO 忽略不存在的teacherID和courseID?
	//resp.Data=BinaryGraphMatchByHungary(req.TeacherCourseRelationShip)
	resp.Data=BinaryGraphMatchByFlow(req.TeacherCourseRelationShip)
}

func BinaryGraphMatchByHungary(G map[string][]string) map[string]string{
	/*
		二分图匹配by匈牙利
		传入图，返回匹配结果
		没测过
	*/
	//TODO string转int，拿到结果之后再转回string会不会更快
	//init
	matchL:=make(map[string]string,len(G))	//左半部匹配情况
	matchR:=make(map[string]string,len(G))	//右半部匹配情况
	//
	var mark map[string]bool
	var dfs func(x string) bool
	dfs=func(x string)bool{
		for _,v:=range G[x]{
			if !mark[v]{
				mark[v]=true
				nowCourseID,ok:=matchR[v]
				if ok==false || dfs(nowCourseID){
					matchR[v]=x
					matchL[x]=v
					return true
				}
			}
		}
		return false
	}
	for x,_:=range G{
		mark=make(map[string]bool,len(G))
		dfs(x)
	}
	return matchL
}
func BinaryGraphMatchByFlow(G map[string][]string) map[string]string{
	/*
		二分图匹配by网络流
		传入图，返回匹配结果
		oj模板题测试maxFlow通过，但是取方案是否正确没测
	*/
	mp:=make(map[string]int,len(G))
	nameL:=make(map[int]string,len(G))
	nameR:=make(map[int]string,len(G))
	n:=len(G)	//左半部[1,n]
	m:=0
	for _,CourseIDArray:=range G{	//右半部CourseID字符串转int
		for _,CourseID:=range CourseIDArray{
			m++	//统计边数
			if _,ok:=mp[CourseID];!ok{
				n++	//统计总点数
				mp[CourseID]=n
				nameR[n]=CourseID
			}
		}
	}
	m+=n
	//build
	f:=&Flow{}
	f.init(n+10,m*2+10,n)
	x:=0
	for TeacherID,CourseIDArray:=range G{
		x++
		nameL[x]=TeacherID
		//st->L
		f.add(f.st,x,1)
		f.add(x,f.st,0)
		//L-R
		for _,CourseID:=range CourseIDArray{
			v:=mp[CourseID]
			f.add(x,v,1)
			f.add(v,x,0)
		}
	}
	for _,v:=range mp{
		//R->ed
		f.add(v,f.ed,1)
		f.add(f.ed,v,0)
	}
	//
	f.run()
	//int转回string
	match:=make(map[string]string,len(G))
	for i:=1;i<=n;i++{
		if f.matchL[i]!=0 {
			match[nameL[i]] = nameR[f.matchL[i]]
		}
	}
	return match
}

const INF = int(1e9)
type Flow struct{
	st,ed int
	//maxFlow int	//最大匹配数
	head,nt,to,w,d []int
	cnt int
	matchL,matchR []int
	q []int
	h,t int
	idx int
}
func (f *Flow) add(x,y,z int){
	f.cnt++;f.nt[f.cnt]=f.head[x];f.head[x]=f.cnt;f.to[f.cnt]=y;f.w[f.cnt]=z
}
func (f *Flow) init(n,m,idx int){
	f.cnt=1
	f.head=make([]int,n)
	f.matchL=make([]int,n)
	f.matchR=make([]int,n)
	f.d=make([]int,n)
	f.q=make([]int,n)

	f.nt=make([]int,m)
	f.to=make([]int,m)
	f.w=make([]int,m)
	f.idx=idx
	f.idx++;f.st=idx
	f.idx++;f.ed=idx
}

func (f *Flow) bfs()bool{
	f.h=0;f.t=0
	f.q[f.t]=f.st;f.t++
	for i:=int(0);i<=f.idx;i++{
		f.d[i]=0
	}
	f.d[f.st]=1
	for f.h<f.t{
		x:=f.q[f.h];f.h++
		for i:=f.head[x];i>0;i=f.nt[i]{
			v:=f.to[i]
			if f.w[i]!=0 && f.d[v]==0{
				f.d[v]=f.d[x]+1
				if v==f.ed{
					return true
				}
				f.q[f.t]=v;f.t++
			}
		}
	}
	return f.d[f.ed]!=0
}
func (f *Flow) dfs(x,flow int)int{
	if x==f.ed{
		return flow
	}
	res:=flow
	for i:=f.head[x];i>0;i=f.nt[i]{
		v:=f.to[i]
		if f.w[i]>0&&f.d[v]==f.d[x]+1{
			k:=f.dfs(v,util.Min(res,f.w[i]))
			f.w[i]-=k
			f.w[i^1]+=k
			res-=k
			if k==0{
				f.d[v]=-1
			}
			if res==0{
				break
			}
		}
	}
	return flow-res
}
func (f *Flow) cal() {
	/*
		遍历残量网络获得方案
		正边是偶数索引，回边是奇数
	*/
	cnt:=f.idx-2	//去掉源点和汇点
	for x:=int(1);x<=cnt;x++{
		for i:=f.head[x];i>0;i=f.nt[i]{
			v:=f.to[i]
			if i%2==1&&f.w[i]>0{	//回边
				f.matchR[x]=v
				f.matchL[v]=x
			}
		}
	}
}
func (f *Flow) run(){
	for f.bfs(){
		f.dfs(f.st,INF)
		//f.maxFlow+=f.dfs(f.st,INF)
	}
	f.cal()
}
