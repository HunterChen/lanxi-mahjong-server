package data

import (
	"gopkg.in/mgo.v2/bson"

	"errors"
	"time"
)

// 代理商用户数据
// 代理商列表定义，在代理商页面填写完用户信息即为代理商，纯绑定了上级代理ID的付费玩家不算作代理商。
type Agent_User struct {
	Gameid   string `bson:"_id"`      // 游戏id
	Phone    string `bson:"Phone"`    // 绑定的手机号码
	Nickname string `bson:"Nickname"` //
	Wechat   string `bson:"Wechat"`   // 微信号
	Auth     string `bson:"Auth"`     // 密码验证码
	Pwd      string `bson:"Pwd"`      // MD5密码
	CIP      uint32 `bson:"CIP"`      // 注册账户时的IP地址
	CTime    uint32 `bson:"CTime"`    // 注册时间为代理商的时间
	JTime    uint32 `bson:"JTime"`    //	绑定时间
	Addr     string `bson:"Addr"`     //物理地址
	LIP      uint32 `bson:"LIP"`      // 最后一次登录IP
	LTime    uint32 `bson:"LTime"`    // 最后一次登录IP
	Lv       uint32                   // 代理等级
	Parent   string `bson:"Parent"`   // 父代理ID
	Charge   uint32 `bson:"Charge"`   //充值金额(单位/分)
	Status   uint32 `bson:"Status"`   // 正常0  锁定1  黑名单2

	Balance uint32 `bson:"Balance"` // 用于提现的余额(单位/分)
}



type Agent_Users []*Agent_User

// 增加充值金额
func (this *Agent_User) IncCharge(charge uint32) error {
	if charge <= 0 {
		return errors.New("charge can not  empty 0")

	}
	_, err := C(_AGENT_USER).UpsertId(this.Gameid, bson.M{"$inc": bson.M{"Charge": charge}})
	return err
}

func (this *Agent_Users) GetByUid(uid string, page int, limit int) (int, error) {
	if page < 1 {
		page = 1
	}
	if limit < LIMIT_MIN {
		limit = LIMIT_MIN
	} else if limit > LIMIT_MAX {
		limit = LIMIT_MAX
	}
	err := C(_AGENT_USER).FindId(uid).Sort("-CTime").Skip((page - 1) * limit).Limit(limit).All(this)
	return 1, err
}
func (this *Agent_Users) GetByPhone(phone string, page int, limit int) (int, error) {
	if page < 1 {
		page = 1
	}
	if limit < LIMIT_MIN {
		limit = LIMIT_MIN
	} else if limit > LIMIT_MAX {
		limit = LIMIT_MAX
	}
	err := C(_AGENT_USER).Find(bson.M{"Phone": phone}).Sort("-CTime").Skip((page - 1) * limit).Limit(limit).All(this)
	return 1, err
}
func (this *Agent_User) Get() error {
	if this.Gameid == "" {
		return errors.New("Gameid number can not empty")
	}
	return C(_AGENT_USER).FindId(this.Gameid).One(this)
}


// 向上获取两级父级
func (this *Agent_Users) GetParent2(self string) {
	for i := uint32(1); i <= 2; i++ {
		user := &Agent_User{Gameid: self}
		if err := user.Get(); err == nil {
			self = user.Parent
			user.Lv = i
			*this = append(*this, user)
		}
	}
}


// 判断指定id是否在自己的下级里
func (this *Agent_User) IsJunior(gameid string) bool {
	list := make([]*Agent_User, 0)
	err := C(_AGENT_USER).Find(bson.M{"Parent": this.Gameid}).All(&list)
	if err != nil || len(list) == 0 {
		return false
	}

	// 判断是否在一级下级里
	for _, v := range list {
		if v.Gameid == gameid {
			return true
		}
	}

	// 判断是否在二级及以上的下级里
	for _, v := range list {
		if this.IsJunior(v.Gameid) {
			return true
		}
	}
	return false
}

// 绑定上级
func (this *Agent_User) Bind(parentUid string) error {
	if parentUid == "" {
		return errors.New("parentUid can not  empty")
	}

	if this.Gameid == "" {
		return errors.New("Gameid can not  empty")
	}

	count, _ := C(_AGENT_USER).FindId(parentUid).Count()
	if count == 0 {
		return errors.New("Uid not exists")
	}

	/*if this.IsJunior(parentUid) {
		return errors.New("can not bind junior")
	}
*/
	this.Parent = parentUid

	this.JTime = uint32(time.Now().Unix())

	bs := bson.M{
		"Parent": this.Parent,
		"JTime":  this.JTime,
	}

	this.CTime = uint32(time.Now().Unix())
	_, err := C(_AGENT_USER).UpsertId(this.Gameid, bs)
	return err
}
