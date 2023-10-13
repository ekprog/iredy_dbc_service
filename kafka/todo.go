package kafka

import (
	"microservice/app/kafka"
	"microservice/pkg/pb/events"
)

var TaskDoneTopic *kafka.KafkaTopic[*events.TaskDoneEvent]

func init() {
	var err error

	TaskDoneTopic, err = kafka.Topic[*events.TaskDoneEvent]("task_done_topic")
	if err != nil {
		panic(err)
	}
}
