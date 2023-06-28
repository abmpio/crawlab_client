package sdk

import "go.mongodb.org/mongo-driver/bson/primitive"

type Spider struct {
	id  primitive.ObjectID
	cmd string
}

func (s *Spider) Id() primitive.ObjectID {
	return s.id
}

func (s *Spider) Cmd() string {
	return s.cmd
}
