package rest

//@Service
type TankService struct {
	Bean
	tankDao *TankDao
}

//初始化方法
func (this *TankService) Init(context *Context) {

	//手动装填本实例的Bean. 这里必须要用中间变量方可。
	b := context.GetBean(this.tankDao)
	if b, ok := b.(*TankDao); ok {
		this.tankDao = b
	}

}
