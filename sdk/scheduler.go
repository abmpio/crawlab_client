package sdk

import (
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	Id          primitive.ObjectID   `json:"_id,omitempty"`
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	SpiderId    primitive.ObjectID   `json:"spider_id,omitempty"`
	Cron        string               `json:"cron,omitempty"`
	EntryId     cron.EntryID         `json:"entry_id,omitempty"`
	Cmd         string               `json:"cmd,omitempty"`
	Param       string               `json:"param,omitempty"`
	Mode        string               `json:"mode,omitempty"`
	NodeIds     []primitive.ObjectID `json:"node_ids,omitempty"`
	Priority    int                  `json:"priority,omitempty"`
	Enabled     bool                 `json:"enabled,omitempty"`
	UserId      primitive.ObjectID   `json:"user_id,omitempty"`
}
