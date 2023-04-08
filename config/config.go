package config

type KafkaConfig struct {
	BootstrapServers string
	Topic            string
	ConsumerGroup    string
}

type Projector struct {
	AutoMigration bool
	Kafka         *KafkaConfig
}

type Command struct {
	AutoMigration bool
	EventStoreDB  *DBConfig
	Kafka         *KafkaConfig
}

type Query struct {
	AutoMigration bool
	EntityStoreDB *DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}
