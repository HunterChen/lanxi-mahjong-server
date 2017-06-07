/**********************************************************
* Author        : Michael
* Email         : dolotech@163.com
* Last modified : 2016-06-11 16:24
* Filename      : postbox.go
* Description   : 邮箱
* *******************************************************/
package request

import (
	"lib/socket"
	"data"

	"interfacer"
	"players"
	"protocol"
	"resource"

	"code.google.com/p/goprotobuf/proto"
)

func init() {

	socket.Regist(&protocol.CPost{}, getAllPost)

	socket.Regist( &protocol.CDelPost{}, delPost)

	socket.Regist(&protocol.CDelReadPost{}, delReadPost)

	socket.Regist( &protocol.CDelAllPost{}, delAllPost)

	socket.Regist(&protocol.COpenAppendix{}, openAppendix)

	socket.Regist(&protocol.CReadPost{}, readPost)

}

func getAllPost(ctos *protocol.CPost, c interfacer.IConn) {
	stoc := &protocol.SPost{}
	data := &data.DataPostbox{Receiver: c.GetUserid()}
	list, err := data.ReadAll()
	if err != nil || len(list) == 0 {
		stoc.Error = proto.Uint32(uint32(protocol.Error_PostboxEmpty))
	} else {
		for _, v := range list {
			p := &protocol.PostBoxData{
				Id:           &v.Id,
				Senderuserid: &v.Sender,
				Content:      &v.Content,
				Title:        &v.Title,
				Appendixname: &v.Appendixname,
				Expire:       &v.Expire,
				Read:         &v.Read,
				Kind:         &v.Kind,
				Draw:         &v.Draw,
			}
			stoc.Data = append(stoc.Data, p)
		}
	}
	c.Send(stoc)

}

func delPost(ctos *protocol.CDelPost, c interfacer.IConn) {
	stoc := &protocol.SDelPost{Postid: ctos.Postid}
	if ctos.GetPostid() > 0 {
		database := &data.DataPostbox{Id: ctos.GetPostid(), Receiver: c.GetUserid()}
		err := database.Delete()
		if err != nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_PostNotExist))
		} else {
			stoc.Error = proto.Uint32(0)
		}
	} else {
		stoc.Error = proto.Uint32(uint32(protocol.Error_PostNotExist))
	}
	c.Send(stoc)
}
func delReadPost(ctos *protocol.CDelReadPost, c interfacer.IConn) {
	stoc := &protocol.SDelReadPost{}
	database := &data.DataPostbox{Receiver: c.GetUserid()}
	err := database.CleanupRead()
	if err != nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_PostNotExist))
	} else {
		stoc.Error = proto.Uint32(0)
	}

	c.Send(stoc)

}
func delAllPost(ctos *protocol.CDelAllPost, c interfacer.IConn) {
	stoc := &protocol.SDelAllPost{}
	database := &data.DataPostbox{Receiver: c.GetUserid()}
	err := database.Cleanup()
	if err != nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_PostNotExist))
	} else {
		stoc.Error = proto.Uint32(0)
	}
	c.Send(stoc)

}
func openAppendix(ctos *protocol.COpenAppendix, c interfacer.IConn) {
	stoc := &protocol.SOpenAppendix{}
	//TODO:优化
	player := players.Get(c.GetUserid())
	//
	if ctos.GetPostid() > 0 {
		database := &data.DataPostbox{Id: ctos.GetPostid(), Receiver: c.GetUserid()}
		widgetList, err := database.OpenAppendix()
		if err != nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_AppendixNotExist))
		} else {
			if len(widgetList) > 0 {
				d := &protocol.PostAppendixData{}
				d.Postid = &database.Id
				d.Name = &database.Appendixname
				stoc.Data = d
				m := make(map[uint32]int32)
				for _, v := range widgetList {
					widget := &protocol.WidgetData{}
					widget.Id = &v.Id
					widget.Count = &v.Count
					d.Widgets = append(d.Widgets, widget)
					m[v.Id] = int32(v.Count)
				}
				resource.ChangeMulti(player, m, data.RESTYPE9)
			} else {
				stoc.Error = proto.Uint32(uint32(protocol.Error_AppendixNotExist))
			}
		}

	} else {
		stoc.Error = proto.Uint32(uint32(protocol.Error_PostNotExist))

	}
	c.Send(stoc)

}

func readPost(ctos *protocol.CReadPost, c interfacer.IConn) {
	stoc := &protocol.SReadPost{}
	stoc.Id = ctos.Id
	if ctos.GetId() > 0 {
		database := &data.DataPostbox{Id: ctos.GetId(), Receiver: c.GetUserid()}
		err := database.ReadPost()
		if err != nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_PostNotExist))
		}
	} else {
		stoc.Error = proto.Uint32(uint32(protocol.Error_PostNotExist))
	}

	c.Send(stoc)
}
