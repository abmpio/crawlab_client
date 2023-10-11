package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	jsonUtil "github.com/abmpio/libx/json"
	stringHelper "github.com/abmpio/libx/str"

	"github.com/go-resty/resty/v2"
)

var (
	spiderClient *SpiderClient
)

type baseClient struct {
	baseUrl string
	token   string
}

type SpiderClient struct {
	baseClient
	IScheduleApi
}

func NewSpiderClient(baseUrl string, token string) *SpiderClient {
	return &SpiderClient{
		baseClient: baseClient{
			baseUrl: stringHelper.EnsureEndWith(baseUrl, "/"),
			token:   token,
		},
	}
}

// get client
func Client() *SpiderClient {
	return spiderClient
}

func GlobalClient(client *SpiderClient) {
	spiderClient = client
}

func (c *SpiderClient) GetSpiderByName(name string) (*Spider, error) {
	if len(name) <= 0 {
		return nil, nil
	}
	conditionValue := map[string]string{
		"key":   "name",
		"op":    "eq",
		"value": name,
	}
	v := url.Values{}
	v.Set("page", "1")
	v.Set("size", "1")
	v.Set("conditions", jsonUtil.ObjectToJson(conditionValue))
	v.Set("stats", "false")

	apiPath := fmt.Sprintf("spiders?%s", v.Encode())
	response, err := c.doGet(apiPath, func(o *requestOptions) {
		o.queryParams = &v
	})
	if err != nil {
		return nil, err
	}
	spiderList := make([]Spider, 0)
	if response.Data != nil {
		err = jsonUtil.ConvertObjectTo(response.Data, &spiderList)
		if err != nil {
			return nil, err
		}
	}
	if len(spiderList) <= 0 {
		return nil, nil
	}
	return &spiderList[0], nil
}

func (c *baseClient) doPost(apiPath string, opts ...func(o *requestOptions)) (*spiderResponse, error) {
	// return c.doRequest(apiPath, "POST", data, false)
	return c.doRequestWithResty(apiPath, resty.MethodPost, opts...)
}

func (c *baseClient) doPut(apiPath string, opts ...func(o *requestOptions)) (*spiderResponse, error) {
	// return c.doRequest(apiPath, "PUT", data, false)
	return c.doRequestWithResty(apiPath, resty.MethodPut, opts...)
}

func (c *baseClient) doDelete(apiPath string, opts ...func(o *requestOptions)) (*spiderResponse, error) {
	// return c.doRequest(apiPath, "DELETE", data, true)
	optList := make([]func(o *requestOptions), 0)
	optList = append(optList, opts...)
	optList = append(optList, func(o *requestOptions) {
		o.ignoreResponse = true
	})
	return c.doRequestWithResty(apiPath, resty.MethodDelete, optList...)
}

func (c *baseClient) doGet(apiPath string, opts ...func(o *requestOptions)) (*spiderResponse, error) {
	// return c.doRequest(apiPath, "GET", nil, false)
	return c.doRequestWithResty(apiPath, resty.MethodGet, opts...)
}

func (c *baseClient) doRequest(apiPath string, httpMethod string, data interface{}, ignoreResponse bool) (*spiderResponse, error) {
	url := fmt.Sprintf("%s%s", c.baseUrl, apiPath)
	var payload io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			err = fmt.Errorf("无效的参数,err:%s", err)
			return nil, err
		}

		payload = strings.NewReader(string(jsonData))
	}

	client := &http.Client{}
	req, err := http.NewRequest(httpMethod, url, payload)
	if err != nil {
		err = fmt.Errorf("向服务器发送post请求时返回异常,url:%s,异常信息:%s", url, err.Error())
		return nil, err
	}
	if len(c.token) > 0 {
		req.Header.Add("Authorization", c.token)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	spiderResponse := &spiderResponse{}
	if !ignoreResponse {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, spiderResponse)
		if err != nil {
			err = fmt.Errorf("向服务器发送请求时接收到的数据不是正确的数据,返回的数据为:%s,url:%s", body, url)
			return nil, err
		}
		if !spiderResponse.IsSuccessful() {
			err = fmt.Errorf("向服务器发送请求返回了错误的结果,返回的数据为:%s,url:%s", body, url)
			return nil, err
		}
	}
	return spiderResponse, nil
}

type requestOptions struct {
	queryParams *url.Values
	bodyValue   interface{}

	ignoreResponse bool
}

func (c *baseClient) doRequestWithResty(apiPath string, httpMethod string, opts ...func(o *requestOptions)) (*spiderResponse, error) {
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
	if requestOptions.bodyValue != nil {
		r.SetBody(requestOptions.bodyValue)
	}
	resp, err := r.Execute(httpMethod, url)
	if err != nil {
		err = fmt.Errorf("向服务器发送post请求时返回异常,url:%s,异常信息:%s", url, err.Error())
		return nil, err
	}
	if !resp.IsSuccess() {
		err = fmt.Errorf("向服务器发送请求返回了错误的结果,url:%s,statusCode:%d,返回的body:%s",
			url,
			resp.StatusCode(),
			resp.String())
		return nil, err
	}
	var spiderRes *spiderResponse
	if !requestOptions.ignoreResponse {
		var ok bool
		spiderRes, ok = resp.Result().(*spiderResponse)
		if !ok {
			err = fmt.Errorf("向服务器发送请求时接收到的数据不是正确的数据,url:%s,返回的数据为:%s",
				resp.Status(),
				url)
			return nil, err
		}
	}
	return spiderRes, nil
}
