package sdk

import "go.mongodb.org/mongo-driver/bson/primitive"

type IScheduleApi interface {
	CreateSchedule(schedule *Schedule) (*Schedule, error)
	DisableScheduler(id primitive.ObjectID) error
}
