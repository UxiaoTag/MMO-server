package mgr

import (
	"MMO-server/model"

	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
)

var (
	singleSpaceManager = singleton.Singleton{}
)

type SpaceManager struct {
	//字典
	dict map[int]*model.Space
}

func (sm *SpaceManager) Init() {
	for k, s := range GetDataManagerInstacne().Spaces {
		sm.dict[k] = model.NewSpace(s)
		logger.SLCInfo("初始化地图:%s", s.Name)
	}
}

func (sm *SpaceManager) GetSpace(spaceId int) *model.Space {
	s, ok := sm.dict[spaceId]
	if ok {
		return s
	}
	return nil
}

func GetSpaceManagerInstance() *SpaceManager {
	result, _ := singleton.GetOrDo[*SpaceManager](&singleSpaceManager, func() (*SpaceManager, error) {
		return &SpaceManager{
			dict: map[int]*model.Space{},
		}, nil
	})

	return result
}
