package sdk

import "go.mongodb.org/mongo-driver/bson/primitive"

type Spider struct {
	Id    primitive.ObjectID `json:"_id"`
	Name  string             `json:"name"`
	Cmd   string             `json:"cmd"`
	Param string             `json:"param"`
	Mode  string             `json:"mode"`
}
