package simply

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewEventStore(db *gorm.DB, ec chan Event) EventStore {
	return &eventStore{
		db:           db,
		eventChannel: ec,
	}
}

type EventStore interface {
	AddAndPublishEvent(userNo uint, event Event) (*EventEntity, error)
	GetAllEvents(aggregateId uuid.UUID) ([]*EventEntity, error)
	GetLastEvent(aggregateId uuid.UUID) (*EventEntity, error)
}

type eventStore struct {
	db           *gorm.DB
	eventChannel chan Event
}

func (evs *eventStore) GetLastEvent(aggregateId uuid.UUID) (event *EventEntity, err error) {
	return event, evs.db.Where("aggregate_id = ?", aggregateId).Last(&event).Error
}

func (evs *eventStore) GetAllEvents(aggregateId uuid.UUID) (events []*EventEntity, err error) {
	return events, evs.db.Where("aggregate_id = ?", aggregateId).Find(&events).Error
}

func (evs *eventStore) AddAndPublishEvent(userNo uint, event Event) (eventEntity *EventEntity, err error) {

	jsonPayload, err := json.Marshal(event)
	if err != nil {
		logger.Println("error marshalling command payload: ", err)
	}

	eventEntity = &EventEntity{
		UserNo:      userNo,
		EventType:   event.Type(),
		AggregateId: event.ID(),
		Payload:     jsonPayload,
	}

	err = evs.db.Create(eventEntity).Error
	if err != nil {
		return nil, err
	}

	// publish event to Projector by channel if it is set
	if evs.eventChannel != nil {
		evs.eventChannel <- event
	}

	return eventEntity, nil
}
