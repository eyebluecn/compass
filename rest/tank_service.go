package rest

import (
	"fmt"
	"time"
	"github.com/json-iterator/go"
)


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

//我们存放在tank服务器的文件资源按照如下的命名规则： /app/blog/yyyy/MM/dd/timestamp
func (this *TankService) getStoreDir() string {

	now := time.Now()
	dateString := now.Format("/2006/01/02")

	timestamp := int64(now.UnixNano() / (1000 * 1000))
	return fmt.Sprintf("/app/compass%s/%d", dateString, timestamp)

}

//去tank服务器请求一个tank的信息回来。
func (this *TankService) HttpFetchUploadToken(filename string, privacy bool, size int64) *TankUploadToken {

	//生成要访问的url
	url := fmt.Sprintf("%s%s",
		CONFIG.TankUrl,
		"/api/alien/fetch/upload/token")
	params := make(map[string]string)
	params["email"] = CONFIG.TankEmail
	params["password"] = CONFIG.TankPassword
	params["filename"] = filename
	params["privacy"] = fmt.Sprintf("%t", privacy)
	params["size"] = fmt.Sprintf("%d", size)
	params["dir"] = this.getStoreDir()

	bytes := HttpPost(url, params)

	//用json的方式输出返回值。
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	result := &WebResult{}
	err := json.Unmarshal(bytes, result)
	if err != nil {
		panic(err)
	}

	data := result.Data
	dataJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}


	tankUploadToken := &TankUploadToken{}
	err = json.Unmarshal(dataJson, tankUploadToken)
	if err != nil {
		panic(err)
	}

	return tankUploadToken;
}
