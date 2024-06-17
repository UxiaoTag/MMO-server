package service

import (
	"MMO-server/mgr"
	"MMO-server/model"
	"math"

	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/vector3"
)

var (
	singleSpaceService = singleton.Singleton{}
)

type SpaceService struct {
}

func GetSpaceServiceInstance() *SpaceService {
	instance, _ := singleton.GetOrDo[*SpaceService](&singleSpaceService, func() (*SpaceService, error) {
		return &SpaceService{}, nil
	})
	return instance
}

func (ss *SpaceService) Start() {
	//初始化地图
	mgr.GetSpaceManagerInstance().Init()
	//位置同步请求
	network.GetMessageRouterInstance().Subscribe("proto.SpaceEntitySyncRequest", network.MessageHandler{Op: ss.spaceEntitySyncRequest})

}

func (ss *SpaceService) GetSpace(id int) *model.Space {
	return mgr.GetSpaceManagerInstance().GetSpace(id)
}

func (ss *SpaceService) spaceEntitySyncRequest(msg network.Msg) {
	sp := msg.Sender.Get("Character")
	if sp == nil {
		return
	} else {
		sp = sp.(*model.Character).Space
	}
	//同步请求信息
	netEntity := msg.Message.(*proto.SpaceEntitySyncRequest).EntitySync.Entity
	netV3 := vector3.NewVector3(float64(netEntity.Position.X), float64(netEntity.Position.Y), float64(netEntity.Position.Z))
	//服务端实际的角色信息
	serEntity := mgr.GetEntityManagerInstance().GetEntity(int(netEntity.Id))
	serV3 := vector3.NewVector3(serEntity.Position().X, serEntity.Position().Y, serEntity.Position().Z)
	//计算距离
	distance := vector3.GetDistance(netV3, serV3)
	//使用服务器计算移速
	netEntity.Speed = int32(serEntity.Speed())
	//计算时间差
	dt := min(serEntity.PositionTime(), 1.0)
	//计算限额
	limit := float64(serEntity.Speed()) * dt * 1.5
	logger.SLCInfo("距离%v，阈值%v，间隔%v", distance, limit, dt)
	if math.IsNaN(distance) || distance > limit {
		resp := &proto.SpaceEntitySyncResponse{
			EntitySync: &proto.NEntitySync{
				Entity: serEntity.EntityData(),
				Force:  true,
			},
		}
		msg.Sender.Send(resp)
		return
	}

	//广播同步信息
	sp.(*model.Space).UpdateEntity(msg.Message.(*proto.SpaceEntitySyncRequest).EntitySync)

}
