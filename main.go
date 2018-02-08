package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"compass/rest"
)

func main() {

	//将运行时参数装填到config中去。
	rest.PrepareConfigs()
	context := rest.NewContext()
	defer context.Destroy()

	http.Handle("/", context.Router)

	dotPort := fmt.Sprintf(":%v", rest.CONFIG.ServerPort)

	info := fmt.Sprintf("App started at http://localhost%v", dotPort)
	rest.LogInfo(info)

	err := http.ListenAndServe(dotPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
