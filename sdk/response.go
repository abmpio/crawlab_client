package sdk

import "errors"

type spiderResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

func (r *spiderResponse) IsSuccessful() bool {
	return r.Message == "success"
}

func (r *spiderResponse) IsError() error {
	if r.Message != "success" {
		return errors.New(r.Error)
	}
	return nil
}
