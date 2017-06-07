package data

import (
	"gopkg.in/mgo.v2/bson"
	"lib/utils"
)

type RebateSettleData struct {
	Money uint32  `bson:"Money"`  // 本周期累计返利金额(单位/分)
	Gameid string `bson:"Gameid"`	// 玩家id
	Time uint32  `bson:"Time"`//  记录创建系统时间戳
}

func (this *RebateSettleData)Inc()  error{
	t:= utils.TimestampSaturday()
	_,err:= C(_REBATE_SETTLE).Upsert(bson.M{"Gameid":this.Gameid,"Time":uint32(t)},bson.M{"$inc":bson.M{"Money":this.Money}})
	return err
}
