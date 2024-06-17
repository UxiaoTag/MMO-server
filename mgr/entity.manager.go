package mgr

import (
	"MMO-server/model"
	"sync"

	"github.com/NumberMan1/common/ns/singleton"
)

var (
	singleEntityManager = singleton.Singleton{}
)

type EntityManager struct {
	index int
	//记录全部Entity对象
	allEntities map[int]*model.Entity
	//记录场景中全部Entity列表
	spaceEntities map[int][]*model.Entity
	mutex         sync.Mutex
}

func GetEntityManagerInstance() *EntityManager {
	result, _ := singleton.GetOrDo[*EntityManager](&singleEntityManager, func() (*EntityManager, error) {
		return &EntityManager{
			index:         0,
			allEntities:   make(map[int]*model.Entity),
			spaceEntities: make(map[int][]*model.Entity),
			mutex:         sync.Mutex{},
		}, nil
	})
	return result
}

func (em *EntityManager) AddEntity(spaceId int, entity *model.Entity) {
	em.mutex.Lock()
	defer em.mutex.Unlock()
	entity.EntityData().Id = int32(em.NewEntityId())
	em.allEntities[int(entity.EntityId())] = entity
	_, ok := em.spaceEntities[spaceId]
	if !ok {
		em.spaceEntities[spaceId] = make([]*model.Entity, 0)
	}
	em.spaceEntities[spaceId] = append(em.spaceEntities[spaceId], entity)
}

func (em *EntityManager) RemoveEntity(spaceId int, entity *model.Entity) {
	em.mutex.Lock()
	delete(em.allEntities, int(entity.EntityId()))
	for i, v := range em.spaceEntities[spaceId] {
		if v == entity {
			em.spaceEntities[spaceId] = append(em.spaceEntities[spaceId][:i], em.spaceEntities[spaceId][i+1:]...)
			break
		}
	}
	em.mutex.Unlock()
}

func (em *EntityManager) GetEntity(entityId int) *model.Entity {
	v, ok := em.allEntities[entityId]
	if ok {
		return v
	} else {
		return nil
	}
}

func (em *EntityManager) NewEntityId() int {
	em.mutex.Lock()
	defer em.mutex.Unlock()
	id := em.index
	em.index += 1
	return id
}
