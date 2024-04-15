package sdk

import (
	"fmt"
	"net/url"

	jsonUtil "github.com/abmpio/libx/json"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-resty/resty/v2"
)

type Spider struct {
	Id    primitive.ObjectID `json:"_id"`
	Name  string             `json:"name"`
	Cmd   string             `json:"cmd"`
	Param string             `json:"param"`
	Mode  string             `json:"mode"`
}

type SpiderClient struct {
	*baseClient
}

func newSpiderClient(client *baseClient) *SpiderClient {
	return &SpiderClient{
		baseClient: client,
	}
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
	response, err := c.doRequestWithResty(apiPath, resty.MethodGet, func(o *requestOptions) {
		o.queryParams = &v
	})
	if err != nil {
		return nil, err
	}
	result, err := c.unmarshalResponseValue(response)
	if err != nil {
		return nil, err
	}
	spiderList := make([]Spider, 0)
	if result.Data != nil {
		err = jsonUtil.ConvertObjectTo(result.Data, &spiderList)
		if err != nil {
			return nil, err
		}
	}
	if len(spiderList) <= 0 {
		return nil, nil
	}
	return &spiderList[0], nil
}
