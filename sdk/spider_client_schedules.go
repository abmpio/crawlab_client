package sdk

import (
	"errors"
	"fmt"

	jsonUtil "github.com/abmpio/libx/json"
	"github.com/abmpio/libx/mapx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IScheduleApi interface {
	CreateSchedule(schedule *Schedule) (*Schedule, error)
	UpdateSchedule(id primitive.ObjectID, data map[string]interface{}) (*Schedule, error)

	DisableScheduler(id primitive.ObjectID) error
	EnableScheduler(id primitive.ObjectID) error
	DeleteScheduler(id primitive.ObjectID) error
	GetScheduler(id primitive.ObjectID) (*Schedule, error)
}

// #region ISchedulerApi members

// create a schedule
func (c *SpiderClient) CreateSchedule(schedule *Schedule) (*Schedule, error) {
	if schedule == nil {
		return nil, errors.New("schedule is nil")
	}
	response, err := c.doPost("schedules", func(o *requestOptions) {
		o.bodyValue = schedule
	})
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
	response, err := c.doPut(apiPath, func(o *requestOptions) {
		o.bodyValue = newData
	})
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
	_, err := c.doPost(apiPath)
	if err != nil {
		return err
	}
	return nil
}

// enable scheduler
func (c *SpiderClient) EnableScheduler(id primitive.ObjectID) error {
	apiPath := fmt.Sprintf("schedules/%s/enable", id.Hex())
	_, err := c.doPost(apiPath)
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
	_, err := c.doDelete(apiPath)
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
