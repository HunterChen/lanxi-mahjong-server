package data

import (
	"time"
	"lib/utils"
)

const (
	RESTYPE1  uint32 = iota //普通场抽成
	RESTYPE2                //普通场打牌
	RESTYPE3                //比赛场
	RESTYPE4                //私人局
	RESTYPE5                //破产
	RESTYPE6                //充值
	RESTYPE7                //签到
	RESTYPE8                //Vip
	RESTYPE9                //邮件
	RESTYPE10               //购买
	RESTYPE11               //兑换
	RESTYPE12               //结算
	RESTYPE13               //活动领奖
	RESTYPE14               //排行榜
	RESTYPE15               //任务
	RESTYPE16               //绑定奖励
	RESTYPE17               //后台
)

type DataResChange struct {
	Userid   string `bson:"Userid"`   //玩家ID
	Kind     uint32 `bson:"Kind"`     //道具、货币种类
	Time     uint32 `bson:"Time"`     //变动时间
	Channel  uint32 `bson:"Channel"`  //获取、扣除渠道
	Residual uint32 `bson:"Residual"` //剩余量
	Count    int32 `bson:"Count"`  // 变数量
}


type DataResChanges []*DataResChange
func (this *DataResChanges) Save(userid string) error {
	var list []interface{}
	for _,v:=range *this{
		v.Time =uint32(time.Now().Unix())
		v.Userid = userid
		list =append(list,utils.Struct2Map(v))
	}
	return C(_RESOURCE_RECORD).Insert(list...)
}
