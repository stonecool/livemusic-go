package chrome

type Category struct {
	name      string
	instances []*Instance
}

func (c *Category) AddInstance(ins *Instance) {
	c.instances = append(c.instances, ins)
}

func (c *Category) GetInstances() []*Instance {
	return c.instances
}
