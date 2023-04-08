package generator

import "easywalk.io/go/simply-cqrs/command"

type EntityGenerator interface {
	CreateEntityAnsSave(events []*command.Event) error
}
