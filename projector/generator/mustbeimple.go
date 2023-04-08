package generator

import "github.com/easywalk/go-simply-cqrs/command"

type EntityGenerator interface {
	CreateEntityAnsSave(events []*command.Event) error
}
