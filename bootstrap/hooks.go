package bootstrap

import (
	"fmt"
	"microservice/app/kafka"
	"reflect"
	"strconv"
	"time"
)

func RunHooks() error {
	KafkaTest()
	return nil
}

type TestModel struct {
	Hello string
}

func KafkaTest() {
	return
	fmt.Println("Hello World from hook!")

	kafkaTestTopic, err := kafka.Topic[*TestModel](reflect.TypeOf(TestModel{}).Name())
	if err != nil {
		panic(err)
	}

	//topics, err := kafka.Topics()
	//if err != nil {
	//	panic(err)
	//}
	//log.Printf("%v", strings.Join(topics, ", "))

	go func() {
		i := 0
		for {
			time.Sleep(time.Second)
			msg := &TestModel{
				Hello: strconv.Itoa(i),
			}

			err = kafkaTestTopic.Produce(msg)
			if err != nil {
				panic(err)
			}
			i++
		}
	}()

	messages, err := kafkaTestTopic.StartPolling()
	if err != nil {
		panic(err)
	}
	for msg := range messages {
		fmt.Printf("Message OK: %v\n", msg.Value)
		err := kafkaTestTopic.CommitOffset(msg)
		if err != nil {
			panic(err)
		}
	}

	select {}
}
