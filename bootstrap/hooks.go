package bootstrap

//type TestModel struct {
//	Random int
//}
//
//func KafkaTest() {
//	return
//	fmt.Println("Hello World from hook!")
//
//	kafkaTestTopic, err := kafka.Topic[*events.TestEvent](reflect.TypeOf(events.TestEvent{}).Name())
//	if err != nil {
//		panic(err)
//	}
//
//	//topics, err := kafka.Topics()
//	//if err != nil {
//	//	panic(err)
//	//}
//	//log.Printf("%v", strings.Join(topics, ", "))
//
//	go func() {
//		i := 0
//		for {
//			time.Sleep(time.Second)
//			msg := &events.TestEvent{
//				Hello: strconv.Itoa(i),
//			}
//
//			err = kafkaTestTopic.Produce(msg)
//			if err != nil {
//				panic(err)
//			}
//			i++
//		}
//	}()
//
//	//messages, err := kafkaTestTopic.StartPolling()
//	//if err != nil {
//	//	panic(err)
//	//}
//	//for msg := range messages {
//	//	fmt.Printf("Message OK: %v\n", msg.Value)
//	//	err := kafkaTestTopic.CommitOffset(msg)
//	//	if err != nil {
//	//		panic(err)
//	//	}
//	//}
//
//	select {}
//}
