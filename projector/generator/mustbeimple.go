package generator

import "github.com/easywalk/simply-go-cqrs/command"

type EntityGenerator interface {
	CreateEntityAnsSave(events []*command.Event) error
}
