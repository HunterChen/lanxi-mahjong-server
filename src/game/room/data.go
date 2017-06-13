package room

import (
	"lib/utils"
)

func NewDeskData(rid, round, expire,  ante, cost uint32, creator,  invitecode string,maizi bool) *DeskData {
	return &DeskData{
		Rid:     rid,
		Ante:    ante,
		Cost:    cost,
		Cid:     creator,
		Expire:  expire,
		Round:   round,
		Code:    invitecode,
		CTime:   uint32(utils.Timestamp()),
		Score:   make(map[string]int32),
		MaiZi:maizi,
	}
}

type DeskData struct {
	Rid     uint32           //房间ID
	Cid     string           //房间创建人
	Expire  uint32           //牌局设定的过期时间
	Code    string           //房间邀请码
	Round   uint32           // 总牌局数
	Ante    uint32           //私人房底分
	Cost    uint32           //创建消耗
	CTime   uint32           //创建时间
	Score   map[string]int32 //私人局用户战绩积分
	MaiZi bool		// 是否买子
}
