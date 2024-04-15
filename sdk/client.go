package sdk

import (
	"errors"
	"fmt"
	"net/url"

	stringHelper "github.com/abmpio/libx/str"
	"github.com/go-resty/resty/v2"
)

var (
	crawlabClient *CrawlabClient
)

type baseClient struct {
	baseUrl string
	token   string
}

type CrawlabClient struct {
	*SpiderClient
	*ScheduleClient
}

func NewClient(baseUrl string, token string) *CrawlabClient {
	baseClient := &baseClient{
		baseUrl: stringHelper.EnsureEndWith(baseUrl, "/"),
		token:   token,
	}
	c := &CrawlabClient{
		SpiderClient:   newSpiderClient(baseClient),
		ScheduleClient: newScheduleClient(baseClient),
	}
	return c
}

// get client
func Client() *CrawlabClient {
	return crawlabClient
}

func GlobalClient(client *CrawlabClient) {
	crawlabClient = client
}

// func (c *baseClient) doPost(apiPath string, opts ...func(o *requestOptions)) (*spiderResponse, error) {
// 	// return c.doRequest(apiPath, "POST", data, false)
// 	return c.doRequestWithResty(apiPath, resty.MethodPost, opts...)
// }

// func (c *baseClient) doPut(apiPath string, opts ...func(o *requestOptions)) (*spiderResponse, error) {
// 	// return c.doRequest(apiPath, "PUT", data, false)
// 	return c.doRequestWithResty(apiPath, resty.MethodPut, opts...)
// }

// func (c *baseClient) doDelete(apiPath string, opts ...func(o *requestOptions)) (*spiderResponse, error) {
// 	// return c.doRequest(apiPath, "DELETE", data, true)
// 	return c.doRequestWithResty(apiPath, resty.MethodDelete, opts...)
// }

// func (c *baseClient) doGet(apiPath string, opts ...func(o *requestOptions)) (*spiderResponse, error) {
// 	// return c.doRequest(apiPath, "GET", nil, false)
// 	return c.doRequestWithResty(apiPath, resty.MethodGet, opts...)
// }

// func (c *baseClient) doRequest(apiPath string, httpMethod string, data interface{}, ignoreResponse bool) (*spiderResponse, error) {
// 	url := fmt.Sprintf("%s%s", c.baseUrl, apiPath)
// 	var payload io.Reader
// 	if data != nil {
// 		jsonData, err := json.Marshal(data)
// 		if err != nil {
// 			err = fmt.Errorf("无效的参数,err:%s", err)
// 			return nil, err
// 		}

// 		payload = strings.NewReader(string(jsonData))
// 	}

// 	client := &http.Client{}
// 	req, err := http.NewRequest(httpMethod, url, payload)
// 	if err != nil {
// 		err = fmt.Errorf("向服务器发送post请求时返回异常,url:%s,异常信息:%s", url, err.Error())
// 		return nil, err
// 	}
// 	if len(c.token) > 0 {
// 		req.Header.Add("Authorization", c.token)
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	spiderResponse := &spiderResponse{}
// 	if !ignoreResponse {
// 		body, err := io.ReadAll(res.Body)
// 		if err != nil {
// 			return nil, err
// 		}
// 		err = json.Unmarshal(body, spiderResponse)
// 		if err != nil {
// 			err = fmt.Errorf("向服务器发送请求时接收到的数据不是正确的数据,返回的数据为:%s,url:%s", body, url)
// 			return nil, err
// 		}
// 		if !spiderResponse.IsSuccessful() {
// 			err = fmt.Errorf("向服务器发送请求返回了错误的结果,返回的数据为:%s,url:%s", body, url)
// 			return nil, err
// 		}
// 	}
// 	return spiderResponse, nil
// }

type requestOptions struct {
	queryParams *url.Values
	bodyValue   interface{}

	ignoreResponse bool
}

func (c *baseClient) doRequestWithResty(apiPath string, httpMethod string, opts ...func(o *requestOptions)) (*resty.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseUrl, apiPath)
	client := resty.New()
	r := client.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json")
	if len(c.token) > 0 {
		r.SetHeader("Authorization", c.token)
	}
	requestOptions := &requestOptions{
		ignoreResponse: false,
	}
	for _, eachOpt := range opts {
		eachOpt(requestOptions)
	}
	// queryParams
	if requestOptions.queryParams != nil {
		r.SetQueryParamsFromValues(*requestOptions.queryParams)
	}
	// ignore response?
	if !requestOptions.ignoreResponse {
		r.SetResult(&spiderResponse{})
	}
	r.SetError(&spiderResponse{})
	if requestOptions.bodyValue != nil {
		r.SetBody(requestOptions.bodyValue)
	}
	resp, err := r.Execute(httpMethod, url)
	return resp, err
}

func (c *baseClient) unmarshalResponseValue(resp *resty.Response) (*spiderResponse, error) {
	if resp.IsSuccess() {
		spiderRes, ok := resp.Result().(*spiderResponse)
		if ok {
			return spiderRes, nil
		}
	} else {
		spiderRes, ok := resp.Error().(*spiderResponse)
		if ok {
			return spiderRes, nil
		}
	}
	body := string(resp.Body())
	return nil, errors.New(body)
}
