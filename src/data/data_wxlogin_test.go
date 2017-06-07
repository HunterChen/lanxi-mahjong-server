//微信登录
package data

import (
	"gopkg.in/mgo.v2"
	"testing"
)

func TestWxLogin(t *testing.T) {
	collection := _WX_LOGIN
	host := "127.0.0.1:2225"
	var se *mgo.Session
	var err error
	se, err = mgo.Dial(host)
	se.DB(_DBNAME).C(collection)

	openid := "o48K-wCLzOgvD7b2_kbllFcNHDmQ"
	d := &WXLogin{
		OpenId: openid,
		//AccessToken: openid,
	}
	err = d.Get()
	//err = d.GetByToken()
	//err = d.Save()
	//err = d.Update()
	t.Log(d, err)

	//u := &User{
	//	Userid: "16007",
	//}
	//err = u.Get()
	//t.Log(u, err)
	roundRecord := &GameOverRoundRecord{}
	err = roundRecord.Push(165086)
	t.Log(err)
	//
	list := GameOverRecords{}
	err = list.Get("112918", 1, 5)
	t.Log(err)
	for _, v := range list {
		t.Logf("%+v", v)
		//for _, v2 := range v.Rounds {
		//	t.Logf("%+v", v2)
		//}
	}
}
