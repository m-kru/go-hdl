package gen

type Container []Generable

func (c *Container) Add(g Generable) {
	*c = append(*c, g)
}

func (c *Container) Get(name string) (Generable, bool) {
	for _, g := range *c {
		if g.Name() == name {
			return g, true
		}
	}
	return nil, false
}
