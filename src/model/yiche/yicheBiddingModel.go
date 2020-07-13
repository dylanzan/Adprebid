package yiche

//以下为Request Body
type YicheRequestBody struct {
	Site   Site   `json:"site"`
	User   User   `json:"user"`
	Imp    []Imp  `json:"imp"`
	Date   string `json:"date"`
	Device Device `json:"device"`
	On     int    `json:"on"`
	Id     string `json:"id"`
}

type Site struct {
	Page string   `json:"page"`
	Cat  []string `json:"cat"`
	Id   int      `json:"id"`
}
type User struct {
	Id string `json:"id"`
}

type Device struct {
	Carrier        string `json:"carrier"`
	Connectiontype int    `json:"connectiontype"`
	Ua             string `json:"ua"`
	Make           string `json:"make"`
	Ip             string `json:"ip"`
	Didmd5         string `json:"didmd5"`
	Osv            string `json:"osv"`
	Devicetype     int    `json:"devicetype"`
	Os             int    `json:"os"`
	Macmd5         string `json:"macmd5"`
}

type Imp struct {
	Brand    int    `json:"brand"`
	Id       string `json:"id"`
	Keyword  []int  `json:"keyword"`
	Pid      int    `json:"pid"`
	Banner   Banner `json:"banner"`
	Cat      string `json:"cat"`
	Area     int    `json:"area"`
	City     int    `json:"city"`
	Model    int    `json:"model"`
	Audience int    `json:"audience"`
	TempId   int    `json:"tempId"`
	CpId     string `json:"cpId"`
	Search   string `json:"search"`
	Pmp      Pmp    `json:"pmp"`
	Temp     []int  `json:"temp"`
	On       int    `json:"on"`
	Tagid    string `json:"tagid"`
	Sku      string `json:"sku"`
}

type Banner struct {
	H int `json:"h"`
	W int `json:"w"`
}

type Pmp struct {
	Deals []Deals `json:"deals"`
}

type Deals struct {
	Id string `json:"id"`
}

//以下为Response Body struct
type YicheResponseBody struct {
	Id      string    `json:"id"`
	Seatbid []Seatbid `json:"seatbid"`
	BidId   string    `json:"bidId"`
}

type Seatbid struct {
	Bid []Bid `json:"bid"`
}

type Bid struct {
	Durl   []string `json:"durl"`
	Id     string   `json:"id"`
	Adid   string   `json:"adid"`
	Curl   string   `json:"curl"`
	DealId string   `json:"dealId"`
}
