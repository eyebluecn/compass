package rest

import (
	"time"
	"fmt"
	"net/http"
	"net/url"
	"bytes"
	"strconv"
	"io/ioutil"
	"encoding/json"
)

const URL_FETCH_UPLOAD_TOKEN = "/api/alien/fetch/upload/token"
const URL_FETCH_DOWNLOAD_TOKEN = "/api/alien/fetch/download/token"
const URL_CONFIRM = "/api/alien/confirm"
const URL_UPLOAD = "/api/alien/upload"
const URL_DOWNLOAD = "/api/alien/download"

//@Service
type TankService struct {
	Bean
	tankDao *TankDao
}

type UploadToken struct {
	Base
	UserUuid   string    `json:"userUuid"`
	FolderUuid string    `json:"folderUuid"`
	MatterUuid string    `json:"matterUuid"`
	ExpireTime time.Time `json:"expireTime"`
	Filename   string    `json:"filename"`
	Privacy    bool      `json:"privacy"`
	Size       int64     `json:"size"`
	Ip         string    `json:"ip"`
}

type WebResultUploadToken struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data *UploadToken `json:"data"`
}

//初始化方法
func (this *TankService) Init(context *Context) {

	//手动装填本实例的Bean. 这里必须要用中间变量方可。
	b := context.GetBean(this.tankDao)
	if b, ok := b.(*TankDao); ok {
		this.tankDao = b
	}

}

//在tank服务器上存放的文件路径规则如下： /app/compass/yyyy/MM/dd/timestamp
func (this *TankService) getStoreDir() string {

	now := time.Now()
	datePath := now.Format("/2006/01/02")
	return fmt.Sprintf("/app/compass"+datePath+"/%d", time.Now().UnixNano())
}

//去远程获取uploadToken
func (this *TankService) HttpFetchUploadToken(filename string, privacy bool, size int64) *Tank {

	//生成client 参数为默认
	client := &http.Client{}

	//生成要访问的url
	requestUrl := CONFIG.TankUrl + URL_FETCH_UPLOAD_TOKEN

	data := url.Values{}
	data.Set("email", CONFIG.TankEmail)
	data.Set("password", CONFIG.TankPassword)
	data.Set("filename", filename)
	data.Set("privacy", strconv.FormatBool(privacy))
	data.Set("size", strconv.FormatInt(size, 10));
	data.Set("dir", this.getStoreDir());

	//提交请求
	reqest, err := http.NewRequest("POST", requestUrl, bytes.NewBufferString(data.Encode()))
	this.PanicError(err)

	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	reqest.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//处理返回结果
	response, err := client.Do(reqest)
	this.PanicError(err)

	bs, err := ioutil.ReadAll(response.Body)
	this.PanicError(err)

	fmt.Println(string(bs))

	webResult := new(WebResultUploadToken)
	err = json.Unmarshal(bs, webResult)
	this.PanicError(err)

	//将结果定位到标准输出 也可以直接打印出来 或者定位到其他地方进行相应的处理
	//stdout := os.Stdout
	//_, err = io.Copy(stdout, response.Body)

	//返回的状态码
	//status := response.StatusCode

	//fmt.Println(status)

	return &Tank{}
}
