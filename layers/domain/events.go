package domain

import "microservice/app/kafka"

var UserScoreChangedTopic *kafka.KafkaTopic[*UserScoreChangedEvent]

type UserScoreChangedEvent struct{}
