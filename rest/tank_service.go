package rest

import (
	"fmt"
	"time"
	"github.com/json-iterator/go"
	"strconv"
)

const (
	URL_TANK_FETCH_UPLOAD_TOKEN   = "/api/alien/fetch/upload/token";
	URL_TANK_FETCH_DOWNLOAD_TOKEN = "/api/alien/fetch/download/token";
	URL_TANK_CONFIRM              = "/api/alien/confirm";
	URL_TANK_UPLOAD               = "/api/alien/upload";
	URL_TANK_DOWNLOAD             = "/api/alien/download";
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

//从一个bytes中去提取一个我们想要的对象出来。这个对象是被WebResult包装的。
func (this *TankService) extractResult(bytes []byte, result interface{}) {

	//用json的方式输出返回值。
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	webResult := &WebResult{}
	err := json.Unmarshal(bytes, webResult)
	if err != nil {
		panic(err)
	}

	data := webResult.Data
	dataJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(dataJson, result)
	if err != nil {
		panic(err)
	}

}

//去tank服务器请求一个uploadToken的信息回来。
func (this *TankService) HttpFetchUploadToken(filename string, privacy bool, size int64, operator *User) *Tank {

	//生成要访问的url
	url := fmt.Sprintf("%s%s",
		CONFIG.TankUrl, URL_TANK_FETCH_UPLOAD_TOKEN)
	params := make(map[string]string)
	params["email"] = CONFIG.TankEmail
	params["password"] = CONFIG.TankPassword
	params["filename"] = filename
	params["privacy"] = fmt.Sprintf("%t", privacy)
	params["size"] = fmt.Sprintf("%d", size)
	params["dir"] = this.getStoreDir()

	bytes := HttpPost(url, params)
	tankUploadToken := &TankUploadToken{}

	this.extractResult(bytes, tankUploadToken)

	tank := &Tank{
		UserUuid: operator.Uuid,
		Name:     tankUploadToken.Filename,
		Size:     tankUploadToken.Size,
		Privacy:  tankUploadToken.Privacy}

	this.tankDao.Create(tank)

	tank.UploadTokenUuid = tankUploadToken.Uuid
	tank.UploadUrl = fmt.Sprintf("%s%s", CONFIG.TankUrl, URL_TANK_UPLOAD)

	return tank;
}

//去tank服务器请求一个uploadToken的信息回来。
func (this *TankService) HttpConfirm(uuid string, matterUuid string) *Tank {

	tank := this.tankDao.CheckByUuid(uuid)
	if (tank.Confirmed) {
		panic("文件已经被确认了，请勿重复操作。")
	}

	//生成要访问的url
	url := fmt.Sprintf("%s%s", CONFIG.TankUrl, URL_TANK_CONFIRM)
	params := make(map[string]string)
	params["email"] = CONFIG.TankEmail
	params["password"] = CONFIG.TankPassword
	params["matterUuid"] = matterUuid

	bytes := HttpPost(url, params)
	tankMatter := &TankMatter{}

	this.extractResult(bytes, tankMatter)

	if (tank.Name != tankMatter.Name) {
		panic("文件名不一致，确认失败。");
	}
	if (tank.Size != tankMatter.Size) {
		panic("文件大小不一致，确认失败。");
	}
	if (tank.Privacy != tankMatter.Privacy) {
		panic("文件公开性不一致，确认失败。");
	}

	tank.MatterUuid = matterUuid
	tank.Confirmed = true
	tank.Url = fmt.Sprintf("%s%s/%s/%s", CONFIG.TankUrl, URL_TANK_DOWNLOAD, matterUuid, tank.Name)

	this.tankDao.Save(tank)

	return tank;
}

//去tank服务器获取一个私有文件的下载url.
func (this *TankService) HttpFetchDownloadUrl(uuid string) string {

	tank := this.tankDao.CheckByUuid(uuid)
	if (!tank.Confirmed) {
		panic("文件尚未确认，无法下载。")
	}

	if (!tank.Privacy) {
		return tank.Url;
	}

	//生成要访问的url
	url := fmt.Sprintf("%s%s", CONFIG.TankUrl, URL_TANK_FETCH_DOWNLOAD_TOKEN)
	params := make(map[string]string)
	params["email"] = CONFIG.TankEmail
	params["password"] = CONFIG.TankPassword
	params["matterUuid"] = tank.MatterUuid
	params["expire"] = strconv.Itoa(86400)

	bytes := HttpPost(url, params)
	tankDownloadToken := &TankDownloadToken{}

	this.extractResult(bytes, tankDownloadToken)

	return fmt.Sprintf("%s?downloadTokenUuid=%s", tank.Url, tankDownloadToken.Uuid);
}
