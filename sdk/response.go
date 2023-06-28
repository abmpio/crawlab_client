package sdk

type spiderResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

func (r *spiderResponse) IsSuccessful() bool {
	return r.Message == "success"
}
