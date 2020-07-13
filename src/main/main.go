package main

import (
	"fmt"
	"net/http"
	"tencentgo/src/controller"
	"tencentgo/src/helpers/config"
	"tencentgo/src/router"
	"time"
)

//项目启动
func main() {

	//初始化配置文件
	config.InitConfig()

	//Controller层启动
	//controller.TencentCtlInit()
	controller.IQiyiCtlInit()

	//启动router
	router := router.ReverseProxyRouter()

	server := http.Server{
		Addr:         "0.0.0.0:"+config.MediaConf.Basic.ListenPort,
		Handler:      router,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	fmt.Println(server.Addr)
	fmt.Println("============================================ AdFastBid ============================================")
	fmt.Println()
	fmt.Println("                                              Start Up                                               ")
	fmt.Println("                    ", time.Now())
	fmt.Println("============================================ AdFastBid ============================================")

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}

}
