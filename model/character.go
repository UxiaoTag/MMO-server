package model

import (
	"MMO-server/database"

	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/vector3"
)

// 角色
type Character struct {
	*Actor
	//当前客户端连接
	Conn network.Connection
	//当前数据库对应的数据库对象
	Data database.DbCharacter
}

func NewCharacter(position, direction vector3.Vector3) *Character {
	return &Character{Actor: NewActor(position, direction)}
}

func CharacterFromDbCharacter(dbCharacter database.DbCharacter) *Character {
	c := &Character{
		Actor: NewActor(vector3.NewVector3(
			float64(dbCharacter.X),
			float64(dbCharacter.Y),
			float64(dbCharacter.Z),
		), vector3.Zero3()),
	}
	c.Id = int(dbCharacter.ID)
	c.Name = dbCharacter.Name
	c.Info = &proto.NCharacter{
		Id:       int32(dbCharacter.ID),
		TypeId:   int32(dbCharacter.JobId),
		EntityId: 0,
		Name:     dbCharacter.Name,
		Level:    int32(dbCharacter.Level),
		Exp:      int64(dbCharacter.Exp),
		SpaceId:  int32(dbCharacter.SpaceId),
		Gold:     dbCharacter.Gold,
		Entity:   nil,
		Hp:       int32(dbCharacter.Hp),
		Mp:       int32(dbCharacter.Mp),
	}
	c.Data = dbCharacter
	c.SetSpeed(3000)
	return c
}
