package youku

type YoukuBidRequest struct {
	Id  string
	Imp []Imp
	//Site
	//App
	//Device
	//User
}

type Imp struct {
	Id       string
	Tagid    string
	Bidfloor float64
	Banner   Banner
	//Video

	//Native
	//Pmp
	//Ext
	Secure int
}

type Banner struct {
	W   int
	H   int
	Pos int
}
