package system

import "errors"

type PlatformAllocationStrategy interface {
	GetSuitablePlatform(platforms []*Platform, train Train) (*Platform, error)
}

type GreedyPlatformAllocationStrategy struct{}

func (g GreedyPlatformAllocationStrategy) GetSuitablePlatform(platforms []*Platform, train Train) (*Platform, error) {
	if len(platforms) == 0 {
		return nil, errors.New("no platform available")
	}

	return platforms[0], nil
}
