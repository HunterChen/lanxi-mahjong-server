//充值交易记录
package data

import (
	"lib/utils"

	"gopkg.in/mgo.v2/bson"
)

// 交易记录
type TradeRecord struct {
	Id      string `bson:"_id"`     //商户订单号(游戏内自定义订单号)
	Transid string `bson:"Transid"` //交易流水号(计费支付平台的交易流水号,微信订单号)
	Userid  string `bson:"Userid"`  //用户在商户应用的唯一标识(userid)
	//--
	Itemid    string `bson:"Itemid"`    //购买商品ID
	Amount    string `bson:"Amount"`    //购买商品数量
	Diamond   uint32 `bson:"Diamond"`   //购买钻石数量
	Money     uint32 `bson:"Money"`     //交易总金额(单位为分)
	Transtime string `bson:"Transtime"` //交易完成时间 yyyy-mm-dd hh24:mi:ss
	Result    int    `bson:"Result"`    //交易结果(0–交易成功,1–交易失败,2-交易中,3-发货中)
	//--iapppay
	Waresid   uint32 `bson:"Waresid"`   //商品编码(平台为应用内需计费商品分配的编码)
	Currency  string `bson:"Currency"`  //货币类型(RMB,CNY)
	Transtype int    `bson:"Transtype"` //交易类型(0–支付交易)
	Feetype   int    `bson:"Feetype"`   //计费方式(表示商品采用的计费方式)
	Paytype   uint32 `bson:"Paytype"`   //支付方式(表示用户采用的支付方式,403-微信支付)
	//--
	Clientip string `bson:"Clientip"` //客户端ip
	Parent   string `bson:"Parent"`   //绑定的父级代理商游戏ID
	Ctime    int64 `bson:"Ctime"`     //本条记录生成unix时间戳
}

// 生成订单id,(时间截+角色id)
func GenCporderid(userid string) string {
	return utils.Base62encode(uint64(utils.TimestampNano())) + userid
}

// 交易结果记录
func (this *TradeRecord) Get() error {
	return C(_TRADE_RECODE).FindId(this.Id).One(this)
}

func (this *TradeRecord) Update() error {
	return C(_TRADE_RECODE).UpdateId(this.Id, utils.Struct2Map(this))
}

func (this *TradeRecord) Save() error {
	return C(_TRADE_RECODE).Insert(this)
}

// 获取某玩家的所有离线订单,用于上线补单
func GetTradeOff(userid string) ([]*TradeRecord, error) {
	list := make([]*TradeRecord, 0)
	err := C(_TRADINGOFFLINE).Find(bson.M{"Userid": userid}).All(&list)
	return list, err
}

// 保存离线订单,用于下次上线补单
func (this *TradeRecord) SaveTradeOff() error {
	return C(_TRADINGOFFLINE).Insert(this)
}

// 补单完成，删除订单
func DelTradeOff(userid string) error {
	return C(_TRADINGOFFLINE).Remove(bson.M{"Userid": userid})

}

// todo 对已绑定上级的玩家充值时累计记录该玩家的父级和祖父级de周期代理返利
func (this *TradeRecord) Settle(gameid string, money uint32) error {
	users := Agent_Users{}
	users.GetParent2(gameid)
	for _, v := range users {
		if v.Lv == 1 {
			rebate := RebateSettleData{Gameid: v.Gameid, Money: uint32(float32(money) * 0.4)}
			rebate.Inc()
		} else if v.Lv == 2 {
			rebate := RebateSettleData{Gameid: v.Gameid, Money: uint32(float32(money) * 0.06)}
			rebate.Inc()
		}
	}

	return nil
}
