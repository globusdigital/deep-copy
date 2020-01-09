package testdata

type Alpha struct {
	B *Beta
	G Gamma
	D *Delta
	E Epsilon
}

type Beta struct {
	ch chan int
}

type Gamma struct{}

type Delta struct{}

type Epsilon struct{}

func (b *Beta) DeepCopy() *Beta {
	cp := &Beta{ch: make(chan int)}

	return cp
}

func (g *Gamma) DeepCopy() Gamma {
	return Gamma{}
}

func (d Delta) DeepCopy() Delta {
	return Delta{}
}

func (e Epsilon) DeepCopy() *Epsilon {
	return &Epsilon{}
}
