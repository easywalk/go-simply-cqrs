package simply

import (
	"encoding/json"
	"github.com/Shopify/sarama"
)

func NewObserver(cfg *KafkaConfig, ec <-chan Event) (Observer, error) {
	// initialize kafka producer
	producer, err := sarama.NewSyncProducer([]string{cfg.BootstrapServers}, nil)
	if err != nil {
		logger.Fatalln("Error initializing kafka", err)
		return nil, err
	}

	// create topic if not exists
	err = createTopic(cfg.BootstrapServers, cfg.Topic)
	if err != nil {
		logger.Fatalln("Error creating topic", err)
		return nil, err
	}

	logger.Println("Kafka Producer initialized")

	var p = &observer{
		producer: producer,
		ec:       ec,
		topic:    cfg.Topic,
	}
	p.Run()

	return p, nil
}

func createTopic(servers string, topic string) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	admin, err := sarama.NewClusterAdmin([]string{servers}, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := admin.Close(); err != nil {
			logger.Fatalf("Error closing admin: %v", err)
		}
	}()

	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     5,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		return err
	}

	return nil
}

type Observer interface {
	Run()
}

type observer struct {
	ec       <-chan Event
	producer sarama.SyncProducer
	topic    string
}

func (p *observer) Run() {
	// event loop
	go func() {
		for {
			select {
			case event := <-p.ec:
				p.publishEvent(&event)
			}
		}
	}()
}

func (p *observer) publishEvent(event *Event) {
	// publish event to kafka
	jsonPayload, err := json.Marshal(event)
	if err != nil {
		logger.Println("error marshalling command payload: ", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("eventType"),
				Value: []byte((*event).Type()),
			},
		},
		Key:   sarama.StringEncoder((*event).ID().String()),
		Value: sarama.StringEncoder(jsonPayload),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		logger.Println("error publishing event: ", err)
	}
}
