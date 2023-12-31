package domain

type BaseResponse struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	ExtCode string `json:"ext_code"`
	ExtInfo string `json:"ext_info"`
}

type LinearBaseResponse struct {
	BaseResponse     `json:",inline"`
	TimeNow          string `json:"time_now"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMs int    `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}
