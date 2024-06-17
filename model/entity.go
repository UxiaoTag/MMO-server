package model

import (
	"time"

	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/vector3"
)

type Entity struct {
	speed      int             //移速
	position   vector3.Vector3 //位置
	direction  vector3.Vector3 //方向
	netObj     *proto.NEntity  //网络对象
	lastUpdate int64           //上次更新时间
}

// PositionTime 距离上次时间间隔/s
func (e *Entity) PositionTime() float64 {
	return float64(time.Now().UnixMilli()-e.lastUpdate) * 0.001
}

func (e *Entity) Speed() int {
	return e.speed
}

func (e *Entity) SetSpeed(speed int) {
	e.speed = speed
}

func (e *Entity) Position() vector3.Vector3 {
	return e.position
}

func (e *Entity) SetPosition(position vector3.Vector3) {
	e.position = position

	//	记录网络位置
	e.netObj.Position = &proto.NVector3{
		X: int32(position.X),
		Y: int32(position.Y),
		Z: int32(position.Z),
	}
	e.lastUpdate = time.Now().UnixMilli()
}

func (e *Entity) Direction() vector3.Vector3 {
	return e.direction
}

func (e *Entity) SetDirection(direction vector3.Vector3) {
	e.direction = direction

	//	记录网络方向
	e.netObj.Direction = &proto.NVector3{
		X: int32(direction.X),
		Y: int32(direction.Y),
		Z: int32(direction.Z),
	}
}

func NewEntity(position, direction vector3.Vector3) *Entity {
	e := &Entity{
		netObj:     &proto.NEntity{},
		lastUpdate: time.Now().UnixMilli(),
	}
	e.SetPosition(position)
	e.SetDirection(direction)
	return e
}

// 返回网络对象id
func (e *Entity) EntityId() int32 {
	return e.netObj.Id
}

// 返回网络对象
func (e *Entity) EntityData() *proto.NEntity {
	return e.netObj
}

func (e *Entity) SetEntityData(entity *proto.NEntity) {
	e.netObj = entity
	e.SetPosition(vector3.NewVector3(
		float64(entity.Position.X),
		float64(entity.Position.Y),
		float64(entity.Position.Z),
	))
	e.SetDirection(vector3.NewVector3(
		float64(entity.Direction.X),
		float64(entity.Direction.Y),
		float64(entity.Direction.Z),
	))

	e.SetSpeed(int(entity.Speed))
}
