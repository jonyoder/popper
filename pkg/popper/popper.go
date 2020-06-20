package popper

type DefaultPopper struct {
	name string
}

func (p *DefaultPopper) Name() string {
	return p.name
}
