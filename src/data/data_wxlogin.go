//微信登录
package data

import (
	"lib/utils"

	"gopkg.in/mgo.v2/bson"
)

// 微信登录数据
type WXLogin struct {
	OpenId       string `bson:"_id"`           //OpenId
	UnionId      string `bson:"unionid"`       //unionid
	AccessToken  string `bson:"access_token"`  //access_token
	RefreshToken string `bson:"refresh_token"` //refresh_token
	ATExpiresIn  int64  `bson:"at_expires_in"` //access_token过期时间(2小时)
	RTExpiresIn  int64  `bson:"rt_expires_in"` //refresh_token过期时间(30天)
	Utime        int64  `bson:"update_time"`   //更新本条记录unix时间戳
	Ctime        int64  `bson:"create_time"`   //生成本条记录unix时间戳
}

// 生成id,(时间截+随机字符串)
func GenWXLoginID() string {
	return utils.Base62encode(uint64(utils.TimestampNano())) +
	utils.Base62encode(uint64(utils.RandUint32()))
}

// 交易结果记录
func (this *WXLogin) Get() error {
	return C(_WX_LOGIN).FindId(this.OpenId).One(this)
}

func (this *WXLogin) GetByToken() error {
	return C(_WX_LOGIN).Find(bson.M{"access_token": this.AccessToken}).One(this)
}

func (this *WXLogin) Update() error {
	return C(_WX_LOGIN).UpdateId(this.OpenId, utils.Struct2Map(this))
}

func (this *WXLogin) Set(access_token string, expires_in, update_time int64) error {
	return C(_WX_LOGIN).UpdateId(this.OpenId,bson.M{"$set":bson.M{
		"AccessToken":access_token,
		"ATExpiresIn":expires_in,
		"Utime":update_time,
	}})
}

func (this *WXLogin) Save() error {
	return C(_WX_LOGIN).Insert(this)
}

// TODO:定时刷新refresh_token
