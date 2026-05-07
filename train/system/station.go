package system

import "errors"

type Station struct {
	platforms                  []*Platform
	platformAllocationStrategy PlatformAllocationStrategy
}

func (s Station) FindPlatformById(platformId string) *Platform {
	for _, platform := range s.platforms {
		if platform.platformId == platformId {
			return platform
		}
	}
	return nil
}

func (s Station) GetTrainForPlatform(platformId string) (*Train, error) {
	platform := s.FindPlatformById(platformId)
	if platform == nil {
		return nil, errors.New("platform not found")
	}

	return platform.Train(), nil
}

func (s Station) GetPlatformForTrain(trainId string) (*Platform, error) {
	for _, platform := range s.platforms {
		train := platform.Train()
		if train == nil {
			continue
		}
		if train.trainId == trainId {
			return platform, nil
		}
	}

	return nil, errors.New("train does not exist")
}

func (s Station) AllocateTrain(train Train) (*Platform, error) {
	platforms := []*Platform{}
	for _, platform := range s.platforms {
		if platform.IsEmpty() {
			platforms = append(platforms, platform)
		}
	}
	platform, err := s.platformAllocationStrategy.GetSuitablePlatform(platforms, train)
	if err != nil {
		return nil, err
	}
	err = platform.AllocateTrain(train)
	if err != nil {
		return nil, err
	}

	return platform, nil
}
