package simply

import (
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "simply: ", log.LstdFlags|log.Lshortfile)
)

func CreateCommander(cfg *Command, ec chan Event) EventStore {
	eventStoreDB, err := InitPostgresOrExit(cfg.EventStoreDB)
	if err != nil {
		logger.Fatalln("Error connecting to database", err)
	}
	eventStore := createEventStoreOrExit(eventStoreDB, ec)
	return eventStore
}

func createEventStoreOrExit(eventStoreDB *gorm.DB, ec chan Event) EventStore {
	// auto migration for EventStore
	if err := eventStoreDB.Migrator().DropTable(&EventEntity{}); err != nil {
		logger.Fatalln("Error dropping table", err)
	}
	logger.Printf("Dropped table %s", EventEntity{}.TableName())

	if err := eventStoreDB.AutoMigrate(&EventEntity{}); err != nil {
		logger.Fatalln("Error auto migrating database", err)
	}
	logger.Printf("Created tables of EventStore : %s", EventEntity{}.TableName())

	// command store
	eventStore := NewEventStore(eventStoreDB, ec)
	logger.Println("Event store initialized")

	return eventStore
}
