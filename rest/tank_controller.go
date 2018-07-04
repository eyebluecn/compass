package rest

import (
	"net/http"
	"strconv"
	"regexp"
	"fmt"
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

//处理一些特殊的接口，比如参数包含在路径中,一般情况下，controller不将参数放在url路径中
func (this *TankController) HandleRoutes(writer http.ResponseWriter, request *http.Request) (func(writer http.ResponseWriter, request *http.Request), bool) {

	path := request.URL.Path

	//匹配 /api/tank/download/{uuid}
	reg := regexp.MustCompile(`^/api/tank/download/([^/]+)$`)
	strs := reg.FindStringSubmatch(path)
	if len(strs) != 2 {
		return nil, false
	} else {
		var f = func(writer http.ResponseWriter, request *http.Request) {
			this.Download(writer, request, strs[1])
		}
		return f, true
	}
}

//注册自己的路由。
func (this *TankController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {

	routeMap := make(map[string]func(writer http.ResponseWriter, request *http.Request))

	//每个Controller需要主动注册自己的路由。
	routeMap["/api/tank/edit"] = this.Wrap(this.Edit, USER_ROLE_ADMIN)
	routeMap["/api/tank/detail"] = this.Wrap(this.Detail, USER_ROLE_USER)
	//获取上传token.
	routeMap["/api/tank/fetch/upload/token"] = this.Wrap(this.FetchUploadToken, USER_ROLE_USER)

	return routeMap
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

//获取上传token
func (this *TankController) FetchUploadToken(writer http.ResponseWriter, request *http.Request) *WebResult {

	operator := this.checkUser(writer, request)

	//文件名。
	filename := request.FormValue("filename")
	if filename == "" {
		panic("文件名必填")
	} else if m, _ := regexp.MatchString(`[<>|*?/\\]`, filename); m {
		panic(fmt.Sprintf(`【%s】不符合要求，文件名中不能包含以下特殊符号：< > | * ? / \`, filename))
	}

	//文件公有或私有
	privacyStr := request.FormValue("privacy")
	var privacy bool
	if privacyStr == "" {
		panic(`文件公有性必填`)
	} else {
		if privacyStr == "true" {
			privacy = true
		} else if privacyStr == "false" {
			privacy = false
		} else {
			panic(`文件公有性不符合规范`)
		}
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

	tank := this.tankService.HttpFetchUploadToken(filename, privacy, size, operator)

	return this.Success(tank)

}

//获取上传token
func (this *TankController) Confirm(writer http.ResponseWriter, request *http.Request) *WebResult {

	//验证参数。
	uuid := request.FormValue("uuid")
	if uuid == "" {
		panic("uuid参数必填")
	}

	//验证参数。
	matterUuid := request.FormValue("matterUuid")
	if matterUuid == "" {
		panic("matterUuid参数必填")
	}

	tank := this.tankService.HttpConfirm(uuid, matterUuid)

	return this.Success(tank)

}

//下载
func (this *TankController) Download(writer http.ResponseWriter, request *http.Request, uuid string) {

	downloadUrl := this.tankService.HttpFetchDownloadUrl(uuid)

	http.Redirect(writer, request, downloadUrl, http.StatusSeeOther)


}
