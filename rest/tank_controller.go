package rest

import (
	"net/http"
	"strconv"
)

type TankController struct {
	BaseController
	tankDao     *TankDao
	tankService *TankService
}

//初始化方法
func (this *TankController) Init(context *Context) {
	this.BaseController.Init(context)

	//手动装填本实例的Bean. 这里必须要用中间变量方可。
	b := context.GetBean(this.tankDao)
	if b, ok := b.(*TankDao); ok {
		this.tankDao = b
	}

	b = context.GetBean(this.tankService)
	if b, ok := b.(*TankService); ok {
		this.tankService = b
	}

}

//注册自己的路由。
func (this *TankController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {

	routeMap := make(map[string]func(writer http.ResponseWriter, request *http.Request))

	//每个Controller需要主动注册自己的路由。
	routeMap["/api/tank/edit"] = this.Wrap(this.Edit, USER_ROLE_ADMINISTRATOR)
	routeMap["/api/tank/fetch/upload/token"] = this.Wrap(this.FetchUploadToken, USER_ROLE_USER)
	routeMap["/api/tank/detail"] = this.Wrap(this.Detail, USER_ROLE_USER)
	return routeMap
}

//获取上传的token
func (this *TankController) FetchUploadToken(writer http.ResponseWriter, request *http.Request) *WebResult {

	//验证参数。
	filename := request.FormValue("filename")
	if filename == "" {
		panic("filename参数必填")
	}

	//私有性
	privacyStr := request.FormValue("privacy")
	privacy := false
	if privacyStr == "true" {
		privacy = true
	}

	//文件大小
	sizeStr := request.FormValue("size")
	var size int64
	if sizeStr == "" {
		panic(`文件大小必填`)
	} else {
		var err error
		size, err = strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			panic(`文件大小不符合规范`)
		}
		if size < 1 {
			panic(`文件大小不符合规范`)
		}
	}

	tank := this.tankService.HttpFetchUploadToken(filename, privacy, size)

	return this.Success(tank)
}

//修改
func (this *TankController) Edit(writer http.ResponseWriter, request *http.Request) *WebResult {

	//验证参数。
	name := request.FormValue("name")
	if name == "" {
		panic("name参数必填")
	}

	return this.Success("hello")
}

//获取详情
func (this *TankController) Detail(writer http.ResponseWriter, request *http.Request) *WebResult {

	uuid := request.FormValue("uuid")
	if uuid == "" {
		panic("uuid参数必填")
	}

	tank := this.tankDao.CheckByUuid(uuid)

	return this.Success(tank)

}
