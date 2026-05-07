package sdk

import (
	"errors"
	"fmt"

	jsonUtil "github.com/abmpio/libx/json"
	"github.com/abmpio/libx/mapx"
	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Schedule struct {
	Id          bson.ObjectID        `json:"_id,omitempty"`
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	SpiderId    bson.ObjectID        `json:"spider_id,omitempty"`
	Cron        string               `json:"cron,omitempty"`
	EntryId     cron.EntryID         `json:"entry_id,omitempty"`
	Cmd         string               `json:"cmd,omitempty"`
	Param       string               `json:"param,omitempty"`
	Mode        string               `json:"mode,omitempty"`
	NodeIds     []bson.ObjectID      `json:"node_ids,omitempty"`
	Priority    int                  `json:"priority,omitempty"`
	Enabled     bool                 `json:"enabled,omitempty"`
	UserId      bson.ObjectID        `json:"user_id,omitempty"`
}

type IScheduleApi interface {
	CreateSchedule(schedule *Schedule) (*Schedule, error)
	UpdateSchedule(id bson.ObjectID, data map[string]interface{}) (*Schedule, error)

	DisableScheduler(id bson.ObjectID) error
	EnableScheduler(id bson.ObjectID) error
	DeleteScheduler(id bson.ObjectID) error
	GetScheduler(id bson.ObjectID) (*Schedule, error)
}

type ScheduleClient struct {
	*baseClient
	IScheduleApi
}

func newScheduleClient(baseClient *baseClient) *ScheduleClient {
	return &ScheduleClient{
		baseClient: baseClient,
	}
}

// #region ISchedulerApi members

// create a schedule
func (c *ScheduleClient) CreateSchedule(schedule *Schedule) (*Schedule, error) {
	if schedule == nil {
		return nil, errors.New("schedule is nil")
	}
	response, err := c.doRequestWithResty("schedules", resty.MethodPost, func(o *requestOptions) {
		o.bodyValue = schedule
	})
	if err != nil {
		return nil, err
	}
	result, err := c.unmarshalResponseValue(response)
	if err != nil {
		return nil, err
	}
	resSchedule := &Schedule{}
	if result.Data != nil {
		err = jsonUtil.ConvertObjectTo(result.Data, resSchedule)
		if err != nil {
			return nil, err
		}
	}
	return resSchedule, nil
}

func (c *ScheduleClient) UpdateSchedule(id bson.ObjectID, data map[string]interface{}) (*Schedule, error) {
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
	response, err := c.doRequestWithResty(apiPath, resty.MethodPut, func(o *requestOptions) {
		o.bodyValue = newData
	})
	if err != nil {
		return nil, err
	}
	result, err := c.unmarshalResponseValue(response)
	if err != nil {
		return nil, err
	}

	resSchedule := &Schedule{}
	if result.Data != nil {
		err = jsonUtil.ConvertObjectTo(result.Data, resSchedule)
		if err != nil {
			return nil, err
		}
	}
	return resSchedule, nil
}

// disable scheduler
func (c *ScheduleClient) DisableScheduler(id bson.ObjectID) error {
	apiPath := fmt.Sprintf("schedules/%s/disable", id.Hex())
	response, err := c.doRequestWithResty(apiPath, resty.MethodPost)
	if err != nil {
		return err
	}
	result, err := c.unmarshalResponseValue(response)
	if err != nil {
		return err
	}
	if err := result.IsError(); err != nil {
		return err
	}
	return nil
}

// enable scheduler
func (c *ScheduleClient) EnableScheduler(id bson.ObjectID) error {
	apiPath := fmt.Sprintf("schedules/%s/enable", id.Hex())
	response, err := c.doRequestWithResty(apiPath, resty.MethodPost)
	if err != nil {
		return err
	}
	result, err := c.unmarshalResponseValue(response)
	if err != nil {
		return err
	}
	if err := result.IsError(); err != nil {
		return err
	}
	return nil
}

func (c *ScheduleClient) DeleteScheduler(id bson.ObjectID) error {
	if id.IsZero() {
		return nil
	}
	apiPath := fmt.Sprintf("schedules/%s", id.Hex())
	response, err := c.doRequestWithResty(apiPath, resty.MethodDelete)
	if err != nil {
		return err
	}
	result, err := c.unmarshalResponseValue(response)
	if err != nil {
		return err
	}
	if err := result.IsError(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) ||
			err.Error() == mongo.ErrNoDocuments.Error() {
			return nil
		}
		return err
	}
	return nil
}

func (c *ScheduleClient) GetScheduler(id bson.ObjectID) (*Schedule, error) {
	if id.IsZero() {
		return nil, nil
	}
	apiPath := fmt.Sprintf("schedules/%s", id.Hex())
	response, err := c.doRequestWithResty(apiPath, resty.MethodGet)
	if err != nil {
		return nil, err
	}
	result, err := c.unmarshalResponseValue(response)
	if err != nil {
		return nil, err
	}

	schedule := &Schedule{}
	if result.Data != nil {
		err = jsonUtil.ConvertObjectTo(result.Data, schedule)
		if err != nil {
			return nil, err
		}
	}
	return schedule, nil
}

// #endregion
