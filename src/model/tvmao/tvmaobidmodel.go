package tvmao

import "time"

//bidRequest Model

type TvMaoBidRequest struct {
	Id         string    `json:"id"`
	ExpireTime time.Time `json:"expireTime"`
	Imp        []Imp     `json:"imp"`
	Device     Device    `json:"device"`
	User       User      `json:"user"`
	At         int       `json:"at"`
}

type Imp struct {
	Id       string  `json:"id"`
	Tagid    string  `json:"tagid"`
	BidFloor float64 `json:"bidfloor"`
	Banner   Banner  `json:"banner"`
	Video    Video   `json:"video"`
	DealId   string  `json:"dealid"`
}

type Banner struct {
	W    int `json:"w"`
	H    int `json:"h"`
	Wmax int `json:"wmax"`
	Wmin int `json:"wmin"`
	Hmax int `json:"hmax"`
	Hmin int `json:"hmin"`
}

type Video struct {
	Mimes       []string `json:"mimes"`
	Linearity   int      `json:"linearity"`
	Minduration int      `json:"minduration"`
	Maxduration int      `json:"maxduration"`
	W           int      `json:"w"`
	H           int      `json:"h"`
}

type Device struct {
	Ip  string `json:"ip"`
	Mac string `json:"mac"`
}

type User struct {
	Id string `json:"id"`
}

//tvmao BidResponse

type TvMaoBidResponse struct {
	Id      string    `json:"id"`
	BidId   string    `json:"bidid"`
	SeatBid []SeatBid `json:"seatbid"`
}

type SeatBid struct {
	Bid Bid `json:"bid"`
}

type Bid struct {
	Id    string  `json:"id"`
	ImpId string  `json:"impid"`
	Price float64 `json:"price"`
	Nurl  string  `json:"nurl"`
	Adm   string  `json:"adm"`
	Ext   Ext     `json:"ext"`
}

type Ext struct {
	Monitorcode []string `json:"monitorcode"`
}
