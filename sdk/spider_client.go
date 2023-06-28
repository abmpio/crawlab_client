package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	jsonUtil "github.com/abmpio/libx/json"
	"github.com/abmpio/libx/mapx"
	stringHelper "github.com/abmpio/libx/str"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	response, err := c.doGet(apiPath)
	if err != nil {
		return nil, err
	}
	spider := &Spider{}
	if response.Data != nil {
		err = jsonUtil.ConvertObjectTo(response.Data, spider)
		if err != nil {
			return nil, err
		}
	}
	return spider, nil
}

// #region ISchedulerApi members

// create a schedule
func (c *SpiderClient) CreateSchedule(schedule *Schedule) (*Schedule, error) {
	if schedule == nil {
		return nil, errors.New("schedule is nil")
	}
	response, err := c.doPost("schedules", schedule)
	if err != nil {
		return nil, err
	}
	resSchedule := &Schedule{}
	if response.Data != nil {
		err = jsonUtil.ConvertObjectTo(response.Data, resSchedule)
		if err != nil {
			return nil, err
		}
	}
	return resSchedule, nil
}

func (c *SpiderClient) UpdateSchedule(id primitive.ObjectID, data map[string]interface{}) (*Schedule, error) {
	schedule, err := c.GetScheduler(id)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, fmt.Errorf("scheduler not exist,id:%s", id.Hex())
	}
	newData := make(map[string]interface{})
	if err := jsonUtil.ConvertObjectTo(schedule, &newData); err != nil {
		return nil, err
	}
	mapx.MergeMaps(data, newData, mapx.MergeConfig{
		OnlyReplaceExist: true,
	})
	apiPath := fmt.Sprintf("schedules/%s", id.Hex())
	response, err := c.doPut(apiPath, newData)
	if err != nil {
		return nil, err
	}
	resSchedule := &Schedule{}
	if response.Data != nil {
		err = jsonUtil.ConvertObjectTo(response.Data, resSchedule)
		if err != nil {
			return nil, err
		}
	}
	return resSchedule, nil
}

// disable scheduler
func (c *SpiderClient) DisableScheduler(id primitive.ObjectID) error {
	apiPath := fmt.Sprintf("schedules/%s/disable", id.Hex())
	_, err := c.doPost(apiPath, nil)
	if err != nil {
		return err
	}
	return nil
}

// enable scheduler
func (c *SpiderClient) EnableScheduler(id primitive.ObjectID) error {
	apiPath := fmt.Sprintf("schedules/%s/enable", id.Hex())
	_, err := c.doPost(apiPath, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *SpiderClient) DeleteScheduler(id primitive.ObjectID) error {
	if id.IsZero() {
		return nil
	}
	apiPath := fmt.Sprintf("schedules/%s", id.Hex())
	_, err := c.doDelete(apiPath, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *SpiderClient) GetScheduler(id primitive.ObjectID) (*Schedule, error) {
	if id.IsZero() {
		return nil, nil
	}
	apiPath := fmt.Sprintf("schedules/%s", id.Hex())
	response, err := c.doGet(apiPath)
	if err != nil {
		return nil, err
	}
	schedule := &Schedule{}
	if response.Data != nil {
		err = jsonUtil.ConvertObjectTo(response.Data, schedule)
		if err != nil {
			return nil, err
		}
	}
	return schedule, nil
}

// #endregion

func (c *baseClient) doPost(apiPath string, data interface{}) (*spiderResponse, error) {
	return c.doRequest(apiPath, "POST", data, false)
}

func (c *baseClient) doPut(apiPath string, data interface{}) (*spiderResponse, error) {
	return c.doRequest(apiPath, "PUT", data, false)
}

func (c *baseClient) doDelete(apiPath string, data interface{}) (*spiderResponse, error) {
	return c.doRequest(apiPath, "DELETE", data, true)
}

func (c *baseClient) doGet(apiPath string) (*spiderResponse, error) {
	return c.doRequest(apiPath, "GET", nil, false)
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
