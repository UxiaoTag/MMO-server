package service

import (
	"MMO-server/database"
	"MMO-server/mgr"
	"strings"
	"unicode/utf8"

	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
)

var (
	singleUserService = singleton.Singleton{}
)

// UserService 玩家服务
// 注册，登录，创建角色，进入游戏
type UserService struct {
}

func GetUserServiceInstance() *UserService {
	instance, _ := singleton.GetOrDo[*UserService](&singleUserService, func() (*UserService, error) {
		return &UserService{}, nil
	})
	return instance
}

func (us *UserService) Start() {
	// network.GetMessageRouterInstance().Subscribe()
}

// 删除角色的请求
func (us *UserService) characterDeleteRequest(msg network.Msg) {
	player := msg.Sender.Get("DbPlayer").(*database.DbPlayer)
	database.OrmDb.Where("id = ?", msg.Message.(*proto.CharacterDeleteRequest).CharacterId).
		Where("player_id=?", player.ID).
		Delete(&database.DbCharacter{})
	//响应客户端
	rsp := &proto.CharacterDeleteResponse{
		Success: true,
		Message: "执行完成",
	}
	msg.Sender.Send(rsp)
}

// 创建角色
func (us *UserService) CharacterCreateRequest(msg network.Msg) {
	logger.SLCInfo("创建角色:%v", msg.Message)
	rsp := &proto.ChracterCreateResponse{
		Success:   false,
		Character: nil,
	}
	player := msg.Sender.Get("DbPlayer")
	if player == nil {
		// 未登录不能创建角色
		logger.SLCInfo("未登录不能创建角色")
		rsp.Message = "未登录不能创建角色"
		msg.Sender.Send(rsp)
		return
	}
	characters := make([]database.DbCharacter, 0)
	tx := database.OrmDb.Where("player_id = ?", player.(*database.DbPlayer).ID).Find(&characters)
	if tx.RowsAffected >= 4 {
		logger.SLCInfo("角色数量最多4个")
		rsp.Message = "角色数量最多4个"
		msg.Sender.Send(rsp)
		return
	}
	msgTemp := msg.Message.(*proto.CharacterCreateRequest)
	nameLen := utf8.RuneCountInString(msgTemp.GetName())
	// 判断角色名是否为空或包含非法字符如空格等
	if nameLen == 0 || strings.ContainsAny(msgTemp.Name, " \t\r\n\\") {
		logger.SLCInfo("创建角色失败，角色名不能为空或包含非法字符如空格等")
		rsp.Message = "判断角色名不能为空或包含非法字符如空格等"
		msg.Sender.Send(rsp)
		return
	}
	//角色名最长7个字
	if nameLen > 7 {
		logger.SLCInfo("创建角色失败，角色名不能超过7个字符")
		rsp.Message = "创建角色失败，角色名不能超过7个字符"
		msg.Sender.Send(rsp)
		return
	}
	//检验角色名是否存在
	tx = database.OrmDb.Where("name = ?", msgTemp.Name).First(&database.DbCharacter{})
	if tx.RowsAffected > 0 {
		logger.SLCInfo("创建角色失败，角色名已存在")
		rsp.Message = "创建角色失败，角色名已存在"
		msg.Sender.Send(rsp)
		return
	}
	dbCharacter := database.NewDbCharacter()
	dbCharacter.Name = msgTemp.Name
	dbCharacter.JobId = int(msgTemp.JobType)
	dbCharacter.SpaceId = 1
	dbCharacter.PlayerId = int(player.(*database.DbPlayer).ID)
	tx = database.OrmDb.Save(dbCharacter)
	if tx.RowsAffected > 0 {
		rsp.Success = true
		rsp.Message = "角色创建成功"
		msg.Sender.Send(rsp)
	}
}

func (us *UserService) userLoginRequest(msg network.Msg) {
	req := msg.Message.(*proto.UserLoginRequest)
	dbPlayer := &database.DbPlayer{}
	result := database.OrmDb.Where("username = ? and password = ?", req.Username, req.Password).First(&dbPlayer)
	rsp := &proto.UserLoginResponse{}
	if result.Error != nil {
		rsp.Success = false
		rsp.Message = result.Error.Error()
		logger.SLCError("DB访问失败,Error:%s", result.Error.Error())
		return
	}
	if result.RowsAffected > 0 {
		rsp.Success = true
		rsp.Message = "登录成功"
		msg.Sender.Set("DbPlayer", dbPlayer) //登录成功，在conn里记录用户信息
	} else {
		rsp.Success = false
		rsp.Message = "用户名或密码错误"
	}
	msg.Sender.Send(rsp)
}

func (us *UserService) gameEnterRequest(msg network.Msg) {
	rsq := msg.Message.(*proto.GameEnterRequest)
	logger.SLCInfo("有玩家进入游戏,角色Id:%d", rsq.CharacterId)
	player := msg.Sender.Get("DbPlayer").(*database.DbPlayer)
	dbRole := database.DbCharacter{}
	database.OrmDb.Where("player_id = ?", player.ID).
		Where("id = ?", rsq.CharacterId).First(&dbRole)
	logger.SLCInfo("dbRole = %v", dbRole)
	character := mgr.GetCharacterManagerInstance().CreateCharacter(dbRole)
	//通知玩家登录成功
	response := &proto.GameEnterResponse{
		Success:   true,
		Entity:    character.EntityData(),
		Character: character.Info,
	}
	msg.Sender.Send(response)
	//将新角色加入到地图
	space := GetSpaceServiceInstance().GetSpace(dbRole.SpaceId) //新手村
	space.CharacterJoin(msg.Sender, character)
}
