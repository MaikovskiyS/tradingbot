package volatilityos

type Data struct {
	Ticker    string `json:"ticker"`
	Side      string `json:"side"`
	Action    string `json:"action"`
	TimeStamp int64
}
