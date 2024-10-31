package chrome

import "github.com/stonecool/livemusic-go/internal/instance"

type Category struct {
	name      string
	instances []*instance.Instance
}

func (c *Category) AddInstance(ins *instance.Instance) {
	c.instances = append(c.instances, ins)
}

func (c *Category) GetInstances() []*instance.Instance {
	return c.instances
}
