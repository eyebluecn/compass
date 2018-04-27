package rest

import (
	"net/http"
	"strings"
	"bytes"
	"mime/multipart"
	"fmt"
	"io/ioutil"
)

//根据一个请求，获取ip.
func GetIpAddress(r *http.Request) string {
	var ipAddress string

	ipAddress = r.RemoteAddr

	if ipAddress != "" {
		ipAddress = strings.Split(ipAddress, ":")[0]
	}

	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			if ip != "" {
				ipAddress = ip
			}
		}
	}
	return ipAddress
}

//直接发起一个POST请求
func HttpPost(uri string, params map[string]string) []byte {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	fmt.Println(uri)
	fmt.Println(params)

	err := writer.Close()
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", uri, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return content
}
