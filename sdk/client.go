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
