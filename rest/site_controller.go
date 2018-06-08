package rest

import (
	"net/http"
	"strconv"
)

type SiteController struct {
	BaseController
	siteDao     *SiteDao
	siteService *SiteService
}

//初始化方法
func (this *SiteController) Init(context *Context) {
	this.BaseController.Init(context)

	//手动装填本实例的Bean. 这里必须要用中间变量方可。
	b := context.GetBean(this.siteDao)
	if b, ok := b.(*SiteDao); ok {
		this.siteDao = b
	}

	c := context.GetBean(this.siteService)
	if c, ok := c.(*SiteService); ok {
		this.siteService = c
	}

}

//注册自己的路由。
func (this *SiteController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {

	routeMap := make(map[string]func(writer http.ResponseWriter, request *http.Request))

	//每个Controller需要主动注册自己的路由。
	routeMap["/api/site/create"] = this.Wrap(this.Create, USER_ROLE_USER)
	routeMap["/api/site/delete"] = this.Wrap(this.Delete, USER_ROLE_USER)
	routeMap["/api/site/edit"] = this.Wrap(this.Edit, USER_ROLE_USER)
	routeMap["/api/site/detail"] = this.Wrap(this.Detail, USER_ROLE_USER)
	routeMap["/api/site/page"] = this.Wrap(this.Page, USER_ROLE_USER)

	return routeMap
}

//创建一个用户
func (this *SiteController) Create(writer http.ResponseWriter, request *http.Request) *WebResult {

	name := request.FormValue("name")
	if name == "" {
		panic("名称必填！")
	}

	url := request.FormValue("url")
	if url == "" {
		panic(`链接必填`)
	}

	user := this.checkUser(writer, request)

	site := &Site{
		UserUuid: user.Uuid,
		Name:     name,
		Url:      url,
		Visible:  true,
	}

	site = this.siteDao.Create(site)

	return this.Success(site)
}

//创建一个用户
func (this *SiteController) Delete(writer http.ResponseWriter, request *http.Request) *WebResult {

	uuid := request.FormValue("uuid")
	if uuid == "" {
		return this.Error("文件的uuid必填")
	}

	site := this.siteDao.CheckByUuid(uuid)

	//判断文件的所属人是否正确
	user := this.checkUser(writer, request)
	if user.Role != USER_ROLE_ADMINISTRATOR && site.UserUuid != user.Uuid {
		return this.Error(RESULT_CODE_UNAUTHORIZED)
	}

	this.siteDao.Delete(site)

	return this.Success("删除成功！")

}

//编辑一个用户的资料。
func (this *SiteController) Edit(writer http.ResponseWriter, request *http.Request) *WebResult {

	uuid := request.FormValue("uuid")
	name := request.FormValue("name")
	if name == "" {
		panic("名称必填！")
	}

	url := request.FormValue("url")
	if url == "" {
		panic(`链接必填`)
	}

	site := this.siteDao.CheckByUuid(uuid)

	site.Name = name
	site.Url = url

	site = this.siteDao.Save(site)

	return this.Success(site)
}

//获取用户详情
func (this *SiteController) Detail(writer http.ResponseWriter, request *http.Request) *WebResult {

	uuid := request.FormValue("uuid")
	if uuid == "" {
		panic("uuid参数必填")
	}

	site := this.siteDao.CheckByUuid(uuid)

	return this.Success(site)

}

//获取用户列表 管理员的权限。
func (this *SiteController) Page(writer http.ResponseWriter, request *http.Request) *WebResult {

	//如果是根目录，那么就传入root.
	pageStr := request.FormValue("page")
	pageSizeStr := request.FormValue("pageSize")

	userUuid := request.FormValue("userUuid")
	name := request.FormValue("name")
	url := request.FormValue("url")
	visibleStr := request.FormValue("visible")

	orderCreateTime := request.FormValue("orderCreateTime")
	orderModifyTime := request.FormValue("orderModifyTime")

	var page int
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}

	pageSize := 200
	if pageSizeStr != "" {
		tmp, err := strconv.Atoi(pageSizeStr)
		if err == nil {
			pageSize = tmp
		}
	}

	sortArray := []OrderPair{
		{
			key:   "create_time",
			value: orderCreateTime,
		},
		{
			key:   "modify_time",
			value: orderModifyTime,
		},
	}

	pager := this.siteDao.Page(page, pageSize, userUuid, name, url, visibleStr, sortArray)

	return this.Success(pager)
}
