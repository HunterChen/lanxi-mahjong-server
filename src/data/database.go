package data

import (
	"gopkg.in/mgo.v2"
	"github.com/golang/glog"
	"time"
	"config"
)

const (
	_DBNAME = "lanxi_db"

	LIMIT_MAX = 200 // 每页最大记录数量
	LIMIT_MIN = 20  // 每页最小记录数量
)

// 所以集合名字
const (
	_USER        = "user"        // 用户集合
	_GEN_USER_ID = "user_id_gen" // 玩家ID自增
	_GEN_ROOM_ID = "room_id_gen" // 房间ID自增
	_TRADINGOFFLINE   = "tradingoffline"
	_GAMEOVER_RECORD  = "gameover_record" // 私人局结算记录用于前端
	_CARD_RECORD      = "card_record" // 打牌记录,以房间ID+当前局数组成字符串作为_id,TODO:_id这样组合会出现覆盖
	_GAMEOVER_PRIVATE = "gameover_private" // 私人
	_RESOURCE_RECORD  = "resource_record" //资源消耗记录
	_TRADE_RECODE = "trade_record"	// 交易记录
	_AGENT_USER    = "agent_user" // 代理商用户集合
	_REBATE_SETTLE = "agent_rebatesettle" // 每周返利结算
	_WX_LOGIN = "wx_login" // 微信登录数据
)



var session *mgo.Session
func C(collection string) *mgo.Collection {
	if session == nil {
		var err error
		session, err = mgo.Dial(config.Opts().Db_addr)

		//defer session.Close()
		if err != nil {
			glog.Fatalln(err)
		}
		//session.SetPoolLimit(9)
		go func() {
			for {
				time.Sleep(time.Second * 9)
				err := session.Ping()
				if err != nil {
					glog.Errorln(err)
					session.Refresh()
				}
			}

		}()
	}
	//	ses := session.Clone()
	//	defer ses.Close()
	return session.DB(_DBNAME).C(collection)
}
