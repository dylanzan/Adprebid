package controller

import (
	"net/http"
	"strings"
	"tencentgo/src/helpers/config"
)

var (
	youkuConfig      config.Youku
	youkuUpstreamMap map[string]upStreamStruct
	youkuDeals       []string
)

type YoukuHandler struct {
}

type youkuTransport struct {
	http.RoundTripper
}

func (this youkuTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {

	return &http.Response{}, nil
}

func (this YoukuHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {

}

func YoukuCtlInit() {

	youkuConfig = config.MediaConf.Youku
	youkuUpstreamMap = make(map[string]upStreamStruct)

	if len(iqiyiConfig.UpstreamAddrs) > 0 {
		for _, v := range iqiyiConfig.UpstreamAddrs {
			upStream := strings.Split(v, "|")
			//golang map 遍历输出无序的，所以加入id
			id := upStream[0]
			usSplit := strings.Split(upStream[1], ",")
			deals := usSplit[1:]
			uss := &upStreamStruct{
				ipAddr: usSplit[0],
				deals:  deals,
			}
			youkuUpstreamMap[id] = *uss
			youkuDeals = append(youkuDeals, deals...)
		}
	}
}
