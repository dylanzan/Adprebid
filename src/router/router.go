/**
* @Author: Dylan
* @Date: 2020/7/2 10:35
 */

package router

import (
	"github.com/xujiajun/gorouter"
	"sync"
	"tencentgo/src/controller"
)

var once sync.Once
var router *gorouter.Router

func ReverseProxyRouter() *gorouter.Router {
	//var tencent = &controller.TencentHandler{}
	var iqiyi = &controller.IQiyiHandler{}
	once.Do(func() {
		router = gorouter.New()

		router.Group("/")
		//router.POST("tencent.htm", tencent.ServeHTTP)
		//router.POST("yiche.htm", nil)
		//router.POST("iqiyi.htm", iqiyi.ServerHTTP)
		router.POST("iqiyi.htm", iqiyi.ServerHTTP204)
	})
	return router
}
