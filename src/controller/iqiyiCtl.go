package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wxnacy/wgo/arrays"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"tencentgo/src/helpers/config"
	iqiyi "tencentgo/src/model/iqiyi"
)

var (
	iqiyiBodyMap = &sync.Map{}

	iqiyiConfig      config.IQiyi
	iqiyiUpstreamMap map[string]upStreamStruct
	iqiyiDeals       []string
)

type IQiyiHandler struct {
}

type iqiyiTransport struct {
	http.RoundTripper
}

var _ http.RoundTripper = &iqiyiTransport{}

func (this iqiyiTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {

	b, err := ioutil.ReadAll(r.Body)

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)

	newRequest := &iqiyi.BidRequest{}

	err = proto.Unmarshal(b, newRequest)

	data, err := proto.Marshal(newRequest)
	body := ioutil.NopCloser(bytes.NewReader(data))
	r.Body = body
	r.ContentLength = int64(len(data))
	r.Header.Set("Content-Length", strconv.Itoa(len(data)))

	resp, err = this.RoundTripper.RoundTrip(r)

	if err != nil {
		return nil, err
	}
	err = r.Body.Close()

	if err != nil {
		return nil, err
	}

	b, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()

	if err != nil {
		return nil, err
	}

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)

	newResponse := &iqiyi.BidResponse{}

	err = proto.Unmarshal(b, newResponse)

	for k, _ := range newRequest.GetImp() {
		dealId := newRequest.GetImp()[k].GetCampaignId()
		if len(newResponse.GetSeatbid()) > 0 && len(newResponse.GetSeatbid()[k].GetBid()) > 0 {
			for kS, vS := range newResponse.GetSeatbid() {
				adid := vS.GetBid()[kS].GetId() // ??
				iqiyiBodyMap.Store(dealId, bodyContent{adid, 0})
			}
		} else {
			bodyMap.Store(dealId, bodyContent{"0", 1})
		}
	}

	fmt.Println("REQREQREQREQ\n" + newRequest.String())
	fmt.Println("RESPRESPRESPRESP\n" + newResponse.String())

	data, err = proto.Marshal(newResponse)
	body = ioutil.NopCloser(bytes.NewReader(data))

	resp.Body = body

	resp.ContentLength = int64(len(data))

	resp.Header.Set("Content-Length", strconv.Itoa(len(data)))
	return resp, nil
}

func (this IQiyiHandler) ServerHTTP204(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)

	newRequest := &iqiyi.BidRequest{}

	err = proto.Unmarshal(b, newRequest)

	res, err := json.Marshal(newRequest)

	fmt.Println(string(res))

	w.WriteHeader(204)
}

func (this IQiyiHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)

	newRequest := &iqiyi.BidRequest{}

	err = proto.Unmarshal(b, newRequest)

	res, err := json.Marshal(newRequest)

	fmt.Println(string(res))

	if err != nil {
		log.Println(err)
	}

	addr := iqiyiConfig.DefaultUpstreamAddr
	//newRequestDealId:=newRequest.GetImp()[0].GetCampaignId()
	for _, v := range newRequest.GetImp() {
		for _, vMap := range iqiyiUpstreamMap {
			for _, vS := range vMap.deals {
				dealId := strconv.Itoa(int(v.GetCampaignId()))
				if strings.Contains(dealId, vS) || strings.Contains(vS, dealId) || vS == dealId {
					addr = vMap.ipAddr
				}
			}
		}
	}

	remote, err := url.Parse("http://" + addr)

	if err != nil {
		panic(err)
	}

	for _, v := range newRequest.GetImp() {
		dealId := v.GetCampaignId()
		bodycontent, ok := iqiyiBodyMap.Load(dealId)

		if ok && arrays.Contains(iqiyiDeals, dealId) != -1 && bodycontent != nil {
			if rand.Intn(iqiyiConfig.TimesBackToSource) > 1 {
				fmt.Printf("%v==>%v\n", v.GetCampaignId(), addr)

				id := newRequest.GetId()
				bidid := v.GetId()
				adid := bodycontent.(bodyContent).body
				price := int32(9000)

				err = proto.Unmarshal(b, newRequest)
				newResponse := &iqiyi.BidResponse{}
				if adid != "0" {

					newResponse = &iqiyi.BidResponse{
						Id: &id,
						Seatbid: []*iqiyi.Seatbid{
							{
								Bid: []*iqiyi.Bid{{
									Id:    &bidid,
									Impid: &bidid,
									Price: &price,
								},
								},
							},
						},
					}
				} else {
					newResponse = &iqiyi.BidResponse{
						Id: &id,
					}
				}

				data, err := proto.Marshal(newResponse)
				if err != nil {
					w.WriteHeader(204)
				}

				w.Write(data)

				fmt.Println("REQREQREQREQ\n" + newRequest.String())
				fmt.Println("RESPRESPRESPRESP\n" + newResponse.String())
				return
			}
		}
	}

	body := ioutil.NopCloser(bytes.NewReader(b))

	r.Body = body
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Transport = &transport{http.DefaultTransport}
	proxy.ServeHTTP(w, r)

}

func IQiyiCtlInit() {
	iqiyiConfig = config.MediaConf.IQiyi
	iqiyiUpstreamMap = make(map[string]upStreamStruct)

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
			iqiyiUpstreamMap[id] = *uss
			iqiyiDeals = append(iqiyiDeals, deals...)
		}
	}
}
