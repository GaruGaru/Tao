package tests

import (
	"github.com/GaruGaru/Tao/providers"
)

type TestEventProvider struct {
	DojoEvents []providers.DojoEvent
}

func (p TestEventProvider) Events(lat float64, lon float64, rng int, sorting string) ([]providers.DojoEvent, error) {
	return p.DojoEvents, nil
}

