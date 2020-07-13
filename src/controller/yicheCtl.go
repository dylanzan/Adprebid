/**
* @Author: Dylan
* @Date: 2020/6/29 10:37
 */

package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"

	"github.com/wxnacy/wgo/arrays"
	"io/ioutil"
	"log"
	"net/http"
	model "tencentgo/src/model/yiche"
	//"net/url"
	"strings"
	"tencentgo/src/helpers/config"
)

var (
	yicheMap         = &sync.Map{}
	yiCheConfig      config.Yiche
	yiCheUpstreamMap map[string]upStreamStruct
	yicheDeals       []string
)

type YicheHandler struct {
}

type yicheTransport struct {
	http.RoundTripper
}

func (this yicheTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {

	requestBody := &model.YicheRequestBody{}

	reqBodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	json.Unmarshal(reqBodyBytes, requestBody)

	newResponse := &model.YicheResponseBody{}

	resp, err = this.RoundTrip(r)

	if err != nil {
		return nil, err
	}

	err = r.Body.Close()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, newResponse)

	if err != nil {
		log.Println(err)
	}

	dealId := requestBody.Imp[0].Pmp.Deals[0].Id
	if len(newResponse.Seatbid) > 0 && len(newResponse.Seatbid[0].Bid) > 0 {

		adid := newResponse.Seatbid[0].Bid[0].Adid

		yicheMap.Store(dealId, bodyContent{adid, 0})

	} else {
		yicheMap.Store(dealId, bodyContent{"0", 1})
	}

	data, err := json.Marshal(newResponse)

	body := ioutil.NopCloser(bytes.NewReader(data))
	resp.Body = body
	resp.ContentLength = int64(len(data))
	resp.Header.Set("Content-Length", strconv.Itoa(len(data)))

	return resp, nil
}

func (this YicheHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {

	var yicheRequestBody = new(model.YicheRequestBody)

	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(jsonBody, yicheRequestBody)

	if err != nil {
		log.Println(err)
	}

	addr := yiCheConfig.DefaultUpstreamAddr

	newRequestDealId := yicheRequestBody.Imp[0].Pmp.Deals[0].Id
	if len(yiCheUpstreamMap) > 0 {
		for _, yumV := range yiCheUpstreamMap {
			for _, impV := range yicheRequestBody.Imp {
				for _, dealsV := range impV.Pmp.Deals {
					if arrays.Contains(yumV, dealsV.Id) != -1 { //此处需要分析deal ID 符不符合，暂定
						addr = yumV.ipAddr
					}
				}
			}
		}
	}

	remote, err := url.Parse("http://" + addr)

	if err != nil {
		panic(err)
	}

	dealid := yicheRequestBody.Imp[0].Pmp.Deals[0].Id

	bodycontent, ok := yicheMap.Load(dealid)

	if ok && arrays.Contains(yicheDeals, dealid) != -1 && bodycontent != nil {
		fmt.Println(newRequestDealId + "==>" + addr)
		id := yicheRequestBody.Id
		bidId := yicheRequestBody.Imp[0].Id
		adid := bodycontent.(bodyContent).body
		//price:=float32(9000)

		err = json.Unmarshal(jsonBody, yicheRequestBody)
		yicheResponse := &model.YicheResponseBody{}

		if adid != "0" {
			yicheResponse = &model.YicheResponseBody{
				Id: id,
				Seatbid: []model.Seatbid{
					{
						Bid: []model.Bid{
							{
								Id:   bidId,
								Adid: adid,
							},
						},
					},
				},
			}
		} else {
			yicheResponse = &model.YicheResponseBody{
				Id: id,
			}
		}

		data, err := json.Marshal(yicheResponse)

		if err != nil {
			w.WriteHeader(204)
		}

		w.Write(data)
		return
	}

	body := ioutil.NopCloser(bytes.NewReader(jsonBody))
	r.Body = body
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Transport = &yicheTransport{http.DefaultTransport}
	proxy.ServeHTTP(w, r)

}

func YicheCtlInit() {

	yiCheConfig = config.MediaConf.Yiche

	yiCheUpstreamMap = make(map[string]upStreamStruct)
	if len(yiCheConfig.UpstreamAddrs) > 0 {
		for _, v := range yiCheConfig.UpstreamAddrs {
			upStream := strings.Split(v, "|")
			//golang map 遍历输出无序的，所以加入id
			id := upStream[0]
			usSplit := strings.Split(upStream[1], ",")
			deals := usSplit[1:]
			uss := &upStreamStruct{
				ipAddr: usSplit[0],
				deals:  deals,
			}
			yiCheUpstreamMap[id] = *uss
			yicheDeals = append(yicheDeals, deals...)
		}
	}

}
