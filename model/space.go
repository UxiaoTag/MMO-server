package model

import (
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/network"
	pt "github.com/NumberMan1/common/summer/protocol/gen/proto"
)

// use for read SpaceDefine.json
type SpaceDefine struct {
	//场景编码
	SID int `json:"SID"`
	//场景名称
	Name string `json:"Name"`
	//资源
	Resource string `json:"Resource"`
	//类型
	Kind string `json:"Kind"`
	//PK
	AllowPK int `json:"AllowPK"`
}

type Space struct {
	Id   int
	Name string
	Def  SpaceDefine
	//场景中全部角色
	characterDict map[int]*Character
	connCharacter map[network.Connection]*Character
}

func NewSpace(def SpaceDefine) *Space {
	return &Space{
		Id:            def.SID,
		Name:          def.Name,
		Def:           def,
		characterDict: make(map[int]*Character),
		connCharacter: make(map[network.Connection]*Character),
	}
}

func (s *Space) CharacterJoin(conn network.Connection, c *Character) {
	logger.SLCInfo("角色进入场景", "角色ID", c.Id, "场景ID", s.Id, "场景名称", s.Name)
	conn.Set("Character", c) //角色存入链接
	c.Space = s
	c.Conn = conn
	s.characterDict[c.Id] = c

	_, ok := s.connCharacter[conn]
	if !ok {
		s.connCharacter[conn] = c
	}
	//新角色广播给其他玩家
	c.Info.Entity = c.EntityData()
	response := &pt.SpaceCharactersEnterResponse{
		SpaceId:       int32(s.Id),
		CharacterList: make([]*pt.NCharacter, 0),
	}
	response.CharacterList = append(response.CharacterList, c.Info)
	for _, v := range s.characterDict {
		if v.Conn != conn {
			v.Conn.Send(response)
		}
	}

	//向新加入的角色发送场景中其他角色的信息
	for _, v := range s.characterDict {
		if v.Conn == conn {
			continue
		}
		response.CharacterList = make([]*pt.NCharacter, 0)
		response.CharacterList = append(response.CharacterList, v.Info)
		conn.Send(response)
	}
}

// CharacterLeave 角色离开地图
// 客户端离线、切换地图
func (s *Space) CharacterLeave(conn network.Connection, c *Character) {
	logger.SLCInfo("角色离开场景", "角色ID", c.EntityId(), "场景ID", s.Id, "场景名称", s.Name)
	delete(s.characterDict, c.Id)
	// delete(s.connCharacter, conn)

	//向场景中的所有客户端广播角色离开的信息。
	response := &pt.SpaceCharacterLeaveResponse{
		EntityId: c.EntityId(),
	}
	for _, v := range s.characterDict {
		v.Conn.Send(response)
	}
}

// 广播更新Entity消息
func (s *Space) UpdateEntity(sync *pt.NEntitySync) {
	logger.SLCInfo("广播更新Entity消息", sync.String())
	for _, v := range s.characterDict {
		// 遍历场景中的所有角色，如果角色ID与同步信息匹配，则更新该角色的数据。
		if v.EntityId() == sync.Entity.Id {
			v.SetEntityData(sync.GetEntity())
			v.Data.X = int(sync.Entity.Position.X)
			v.Data.Y = int(sync.Entity.Position.Y)
			v.Data.Z = int(sync.Entity.Position.Z)
		} else {
			// 如果不匹配，则向其他角色发送更新信息。
			respose := &pt.SpaceEntitySyncResponse{
				EntitySync: sync,
			}
			v.Conn.Send(respose)
		}
	}
}
