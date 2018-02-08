package rest

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
)

//上下文，管理数据库连接，管理所有路由请求，管理所有的单例component.
type Context struct {
	//数据库连接
	DB *gorm.DB
	//处理所有路由请求
	Router *Router
	//各类的Bean Map。这里面是包含ControllerMap中所有元素
	BeanMap map[string]IBean
	//只包含了Controller的map
	ControllerMap map[string]IController
}

func (this *Context) OpenDb() {

	var err error = nil
	this.DB, err = gorm.Open("mysql", CONFIG.MysqlUrl)

	//是否打开sql日志
	this.DB.LogMode(false)
	if err != nil {
		panic("failed to connect mysql database")
	}
}

func (this *Context) CloseDb() {

	if this.DB != nil {
		this.DB.Close()
	}
}

//构造方法
func NewContext() *Context {

	context := &Context{}

	//处理数据库连接的开关。
	context.OpenDb()

	//初始化Map
	context.BeanMap = make(map[string]IBean)
	context.ControllerMap = make(map[string]IController)

	//注册各类Beans.在这个方法里面顺便把Controller装入ControllerMap中去。
	context.registerBeans()

	//初始化每个bean.
	context.initBeans()

	//初始化Router. 这个方法要在Bean注册好了之后才能。
	context.Router = NewRouter(context)

	return context

}

//注册一个Bean
func (this *Context) registerBean(bean IBean) {

	typeOf := reflect.TypeOf(bean)
	typeName := typeOf.String()

	if element, ok := bean.(IBean); ok {

		err := fmt.Sprintf("【%s】已经被注册了，跳过。", typeName)
		if _, ok := this.BeanMap[typeName]; ok {
			LogError(fmt.Sprintf(err))
		} else {
			this.BeanMap[typeName] = element

			//看看是不是controller类型，如果是，那么单独放在ControllerMap中。
			if controller, ok1 := bean.(IController); ok1 {
				this.ControllerMap[typeName] = controller
			}

		}

	} else {
		err := fmt.Sprintf("注册的【%s】不是Bean类型。", typeName)
		panic(err)
	}

}

//注册各个Beans
func (this *Context) registerBeans() {

	//preference
	this.registerBean(new(PreferenceController))
	this.registerBean(new(PreferenceDao))
	this.registerBean(new(PreferenceService))

	//session
	this.registerBean(new(SessionDao))

	//user
	this.registerBean(new(UserController))
	this.registerBean(new(UserDao))

}

//从Map中获取某个Bean.
func (this *Context) GetBean(bean IBean) IBean {

	typeOf := reflect.TypeOf(bean)
	typeName := typeOf.String()

	if val, ok := this.BeanMap[typeName]; ok {
		return val
	} else {
		err := fmt.Sprintf("【%s】没有注册。", typeName)
		panic(err)
	}
}

//初始化每个Bean
func (this *Context) initBeans() {

	for _, bean := range this.BeanMap {
		bean.Init(this)
	}

}

//销毁的方法
func (this *Context) Destroy() {
	this.CloseDb()

}
