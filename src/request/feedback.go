/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-06-11 16:18
 * Filename      : config.go
 * Description   : 反馈
 * *******************************************************/
package request

import (
	"lib/socket"
	"data"

	"interfacer"
	"players"
	"protocol"
	"time"

	"code.google.com/p/goprotobuf/proto"
)

func init() {
	socket.Regist(&protocol.CFeedback{}, fbconfig)
	socket.Regist(&protocol.CNotice{}, notice)
	socket.Regist(&protocol.CActivity{}, getAllActivity)
	socket.Regist(&protocol.CGetActivityRewards{}, getActivityRewards)
	socket.Regist(&protocol.CWechatShare{}, wechatShare)
	socket.Regist(&protocol.CGetCurrency{}, getCurrency)
}

func getCurrency(ctos *protocol.CGetCurrency, c interfacer.IConn) {
	//TODO:优化
	player := players.Get(c.GetUserid())
	//
	stoc := &protocol.SGetCurrency{
		Roomcard: proto.Uint32(player.GetRoomCard()),
	}
	c.Send(stoc)
}

func fbconfig(ctos *protocol.CFeedback, c interfacer.IConn) {
	stoc := &protocol.SFeedback{Error: proto.Uint32(0)}
	database := &data.DataFeedback{
		Userid:     c.GetUserid(),
		Createtime: uint32(time.Now().Unix()),
		Kind:       ctos.GetKind(),
		Content:    ctos.GetContent(),
	}
	if err := database.Save(); err != nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_FeedfackError))
	}
	c.Send(stoc)
}

func notice(ctos *protocol.CNotice, c interfacer.IConn) {
	stoc := &protocol.SNotice{}

	database := &data.DataNotice{}
	list := database.GetList()
	if len(list) == 0 {
		stoc.Error = proto.Uint32(uint32(protocol.Error_NoticeListEnpty))
	}
	for _, v := range list {
		notice := &protocol.Notice{}
		notice.Id = &v.Id
		notice.Type = &v.Type
		notice.Title = &v.Title
		notice.Content = &v.Content
		notice.Time = &v.CTime
		stoc.List = append(stoc.List, notice)
	}
	c.Send(stoc)
}

func getAllActivity(ctos *protocol.CActivity, c interfacer.IConn) {
	stoc := &protocol.SActivity{}
	c.Send(stoc)
}

func getActivityRewards(ctos *protocol.CGetActivityRewards, c interfacer.IConn) {
	stoc := &protocol.SGetActivityRewards{}
	c.Send(stoc)
}

func wechatShare(ctos *protocol.CWechatShare, c interfacer.IConn) {
	stoc := &protocol.SUpdateActivity{}
	c.Send(stoc)
}
