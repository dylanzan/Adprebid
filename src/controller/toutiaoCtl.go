package controller

import (
	"bytes"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/wxnacy/wgo/arrays"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"tencentgo/src/helpers/config"
	toutiao "tencentgo/src/model/toutiao"
)

var (
	toutiaoMap = &sync.Map{}

	toutiaoConfig      config.Toutiao
	toutiaoUpstreamMap map[string]upStreamStruct
	toutiaoDeals       []string
)

type ToutiaoHandler struct {
}

type toutiaoTransport struct {
	http.RoundTripper
}

func (this toutiaoTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {

	b, err := ioutil.ReadAll(r.Body)
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	newRequest := &toutiao.BidRequest{}

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

	b, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()

	if err != nil {
		return nil, err
	}

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)

	newResponse := &toutiao.BidResponse{}

	err = proto.Unmarshal(b, newResponse)

	for k, v := range newRequest.GetAdslots() {
		for _, vD := range v.GetPmp().Deals {
			dealId := vD.GetNewId()
			if len(newResponse.GetSeatbids()) > 0 && len(newResponse.GetSeatbids()[k].Ads[k].GetDealid()) > 0 {
				for _, vS := range newResponse.GetSeatbids() {
					adid := vS.Ads[k].GetId()
					toutiaoMap.Store(dealId, bodyContent{adid, 0})
				}
			} else {
				toutiaoMap.Store(dealId, bodyContent{"0", 1})
			}
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

func (this ToutiaoHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)

	newRequest := &toutiao.BidRequest{}

	err = proto.Unmarshal(b, newRequest)

	if err != nil {
		panic(err)
	}

	addr := toutiaoConfig.DefaultUpstreamAddr

	for _, v := range newRequest.GetAdslots() {
		for _, vMap := range toutiaoUpstreamMap {
			for _, vS := range vMap.deals {
				for _, vD := range v.GetPmp().GetDeals() {
					dealId := strconv.FormatInt(vD.GetNewId(), 10)
					if strings.Contains(dealId, vS) || strings.Contains(vS, dealId) {
						addr = vMap.ipAddr
					}
				}
			}
		}
	}

	remote, err := url.Parse("http://" + addr)

	for _, v := range newRequest.GetAdslots() {
		for _, vD := range v.GetPmp().GetDeals() {
			dealID := vD.GetNewId()
			bodycontent, ok := toutiaoMap.Load(dealID)

			if ok && arrays.Contains(toutiaoDeals, dealID) != -1 && bodycontent != nil {
				if rand.Intn(toutiaoConfig.TimesBackToSource) > 1 {
					fmt.Printf("%v==>%v\n", dealID, addr)

					id := newRequest.GetRequestId()
					//bidid := v.GetId()
					adid := bodycontent.(bodyContent).body
					price := uint32(9000)

					err = proto.Unmarshal(b, newRequest)
					newResponse := &toutiao.BidResponse{}

					if adid != "0" {
						newResponse = &toutiao.BidResponse{
							RequestId: &id,
							Seatbids: []*toutiao.SeatBid{
								{
									Ads: []*toutiao.Bid{
										{
											Id:    &id,
											Price: &price,
										},
									},
								},
							},
						}
					} else {
						newResponse = &toutiao.BidResponse{
							RequestId: &id,
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
	}
	body := ioutil.NopCloser(bytes.NewReader(b))

	r.Body = body
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Transport = &transport{http.DefaultTransport}
	proxy.ServeHTTP(w, r)

}

func ToutiaoCtlInit() {

	toutiaoConfig = config.MediaConf.Toutiao
	toutiaoUpstreamMap = make(map[string]upStreamStruct)

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
			toutiaoUpstreamMap[id] = *uss
			toutiaoDeals = append(toutiaoDeals, deals...)
		}
	}

}
