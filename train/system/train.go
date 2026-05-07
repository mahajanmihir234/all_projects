package system

import (
	"errors"
	"sync"
)

type Train struct {
	trainId string
}

type Platform struct {
	platformId string
	mutex      sync.Mutex
	train      *Train
}

func (p *Platform) IsEmpty() bool {
	return p.train == nil
}

func (p *Platform) Train() *Train {
	return p.train
}

func (p *Platform) AllocateTrain(train Train) error {
	if p.train != nil {
		return errors.New("platform has no space")
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.train = &train
	return nil
}
