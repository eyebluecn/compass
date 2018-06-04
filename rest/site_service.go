package rest

//@Service
type SiteService struct {
	Bean
	siteDao *SiteDao
}

//初始化方法
func (this *SiteService) Init(context *Context) {

	//手动装填本实例的Bean. 这里必须要用中间变量方可。
	b := context.GetBean(this.siteDao)
	if b, ok := b.(*SiteDao); ok {
		this.siteDao = b
	}

}
