package simply

type EntityGenerator interface {
	CreateEntityAnsSave(events []*EventEntity) error
}
