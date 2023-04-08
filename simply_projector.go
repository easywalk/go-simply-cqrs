package simply

import (
	"context"
	"github.com/Shopify/sarama"
)

func CreateProjector(cfg *Projector, ec <-chan Event) Observer {
	p, err := NewObserver(cfg.Kafka, ec)
	if err != nil {
		logger.Fatalln("Error initializing projector", err)
	}
	return p
}

func StartEntityGenerator(evs EventStore, eg EntityGenerator, cfg *KafkaConfig) {
	kCfg := sarama.NewConfig()
	kCfg.Producer.Return.Successes = true
	conString := []string{cfg.BootstrapServers}
	conn, err := sarama.NewClient(conString, kCfg)
	if err != nil {
		logger.Fatalf("Error creating client: %v", err)
	}

	consumer, err := sarama.NewConsumerGroupFromClient(cfg.ConsumerGroup, conn)
	if err != nil {
		logger.Fatalf("Error creating consumer group client: %v", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Fatalf("Error closing consumer: %v", err)
		}
	}()

	ctx := context.Background()

	p := &transformer{
		Evs: evs,
		Eg:  eg,
	}

	for {
		err := consumer.Consume(ctx, []string{cfg.Topic}, p)
		if err != nil {
			logger.Fatalf("Error from consumer: %v", err)
		}
	}
}
