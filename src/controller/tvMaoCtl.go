package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"tencentgo/src/model/tvmao"
)

var (
	tvmaoMap         = &sync.Map{}
	tvMaoConfig      config.TvMao
	tvMaoUpstreamMap map[string]upStreamStruct
	tvMaoDeals       []string
)

type TvMaoHandler struct {
}

type tvMaoTransport struct {
	http.RoundTripper
}

func (this tvMaoTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {

	b, err := ioutil.ReadAll(r.Body)
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)

	newRequest := tvmao.TvMaoBidRequest{}

	err = json.Unmarshal(b, newRequest)

	data, err := json.Marshal(newRequest)

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

	newResponse := &tvmao.TvMaoBidResponse{}

	err = json.Unmarshal(b, newResponse)

	for k, v := range newRequest.Imp {
		dealId := v.DealId
		if len(newResponse.SeatBid) > 0 && len(newResponse.SeatBid[k].Bid.Id) > 0 {
			adid := v.Id
			tvmaoMap.Store(dealId, bodyContent{adid, 0})
		} else {
			tvmaoMap.Store(dealId, bodyContent{"0", 1})
		}
	}

	fmt.Printf("REQREQREQREQ\n %+v", newRequest)
	fmt.Println("RESPRESPRESPRESP\n %+v", newResponse)

	data, err = json.Marshal(newResponse)

	body = ioutil.NopCloser(bytes.NewReader(data))

	resp.Body = body

	resp.ContentLength = int64(len(data))

	resp.Header.Set("Content-Length", strconv.Itoa(len(data)))

	return resp, nil
}

func (this TvMaoHandler) ServerHttp(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)

	newRequest := &tvmao.TvMaoBidRequest{}

	err = json.Unmarshal(b, newRequest)

	if err != nil {
		panic(err)
	}

	addr := tvMaoConfig.DefaultUpstreamAddr

	for _, v := range newRequest.Imp {
		dealId := v.DealId
		for _, vM := range tvMaoUpstreamMap {
			for _, vd := range vM.deals {
				if strings.Contains(dealId, vd) || strings.Contains(vd, dealId) {
					addr = vM.ipAddr
				}
			}
		}
	}

	remote, err := url.Parse("http://" + addr)

	for _, v := range newRequest.Imp {
		dealId := v.DealId
		bodycontent, ok := tvmaoMap.Load(dealId)

		if ok && arrays.Contains(toutiaoDeals, dealId) != -1 && bodycontent != nil {

			if rand.Intn(tvMaoConfig.TimesBackToSource) > 1 {
				fmt.Println("%v==>%v\n", dealId, addr)

				id := newRequest.Id
				bidid := v.Id

				adid := bodycontent.(bodyContent).body

				price := float64(9000)
				//extid:="ssp"+adid

				err = json.Unmarshal(b, newRequest)

				newResponse := &tvmao.TvMaoBidResponse{}
				if adid != "0" {
					newResponse = &tvmao.TvMaoBidResponse{
						Id:    id,
						BidId: bidid,
						SeatBid: []tvmao.SeatBid{
							{
								Bid: tvmao.Bid{
									Id:    bidid,
									ImpId: bidid,
									Price: price,
									//Nurl:  "",
									//Adm:   "",
									//Ext:   tvmao.Ext{Monitorcode:extid},
								},
							},
						},
					}
				} else {
					newResponse = &tvmao.TvMaoBidResponse{
						Id: id,
					}
				}
				data, err := json.Marshal(newResponse)
				if err != nil {
					w.WriteHeader(204)
				}

				w.Write(data)

				fmt.Printf("REQREQREQREQ\n %+v", newRequest)
				fmt.Println("RESPRESPRESPRESP\n %+v", newResponse)
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

func TvMaoCtlInit() {

	tvMaoConfig = config.MediaConf.TvMao
	tvMaoUpstreamMap = make(map[string]upStreamStruct)

	if len(tvMaoConfig.UpstreamAddrs) > 0 {
		for _, v := range tvMaoConfig.UpstreamAddrs {
			upStream := strings.Split(v, "|")
			//golang map 遍历输出无序的，所以加入id
			id := upStream[0]
			usSplit := strings.Split(upStream[1], ",")
			deals := usSplit[1:]
			uss := &upStreamStruct{
				ipAddr: usSplit[0],
				deals:  deals,
			}
			tvMaoUpstreamMap[id] = *uss
			tvMaoDeals = append(toutiaoDeals, deals...)
		}
	}
}
