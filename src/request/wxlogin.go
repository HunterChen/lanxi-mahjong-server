/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-01-23 10:22
 * Filename      : login.go
 * Description   : 微信登录
 * *******************************************************/
package request

import (
	//"lib/socket"
	"lib/utils"
	"data"

	"interfacer"
	"protocol"
	"time"
	"wxapi"

	"github.com/golang/glog"
	"code.google.com/p/goprotobuf/proto"
)

var WxLogin *wxapi.Wxapi //微信登录

func WxLoginInit() {
	cfg := &wxapi.WxapiConfig{
		AppId:           "wx1c3fa8244c190785",
		AppSecret:       "e25e0b8052a6cbe89be8921b96557728",
		AccessUrl:       "https://api.weixin.qq.com/sns/oauth2/access_token",
		RefreshUrl:      "https://api.weixin.qq.com/sns/oauth2/refresh_token",
		UserinfoUrl:     "https://api.weixin.qq.com/sns/userinfo",
		VerifyAccessUrl: "https://api.weixin.qq.com/sns/auth",
	}
	wx, err := wxapi.NewWxapi(cfg)
	if err != nil {
		panic(err)
	}
	WxLogin = wx
}

//func init() {
//	p := protocol.CWechatLogin{}
//	socket.Regist(p.GetCode(), p, weixinLogin)
//}

type wxLoginData struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	//---
	OpenId      string   `json:"openid"`
	Nickname    string   `json:"nickname"`
	Sex         int      `json:"sex"`
	Province    string   `json:"province"`
	City        string   `json:"city"`
	Country     string   `json:"country"`
	HeadImagUrl string   `json:"headimgurl"`
	Privilege   []string `json:"privilege"`
	UnionId     string   `json:"unionid"`
}

func loginData(wxdata wxLoginData, userinfo wxapi.UserInfoResult) wxLoginData {
	wxdata.OpenId      = userinfo.OpenId     
	wxdata.Nickname    = userinfo.Nickname   
	wxdata.Sex         = userinfo.Sex        
	wxdata.Province    = userinfo.Province   
	wxdata.City        = userinfo.City       
	wxdata.Country     = userinfo.Country    
	wxdata.HeadImagUrl = userinfo.HeadImagUrl
	wxdata.Privilege   = userinfo.Privilege  
	wxdata.UnionId     = userinfo.UnionId    
	return wxdata
}

//code登录
func loginByCode(code string) (wxLoginData, error) {
	wxdata := wxLoginData{}
	var access_token string
	var refresh_token string
	var expires_in int
	//获取access_token
	accessResult, err := WxLogin.Auth(code)
	if err != nil {
		return wxdata, err
	}
	access_token = accessResult.AccessToken
	refresh_token = accessResult.RefreshToken
	expires_in = accessResult.ExpiresIn
	//刷新access_token
	if expires_in <= 0 {
		refreshResult, err := WxLogin.Refresh(refresh_token)
		if err != nil {
			return wxdata, err
		}
		access_token = refreshResult.AccessToken
		refresh_token = refreshResult.RefreshToken
		expires_in = refreshResult.ExpiresIn
	}
	//获取个人信息
	userinfoResult, err := WxLogin.UserInfo(accessResult.OpenId, access_token)
	if err != nil {
		return wxdata, err
	}
	var now int64 = utils.Timestamp()
	var aTime int64 = now + int64(expires_in) //2小时过期
	var rTime int64 = now + 29 * 86400 //30天过期
	//更新数据
	d := &data.WXLogin{
		OpenId: userinfoResult.OpenId,
	}
	err = d.Get()
	if err != nil {
		d.Ctime = now
	}
	d.UnionId = userinfoResult.UnionId
	d.AccessToken  = access_token
	d.RefreshToken = refresh_token
	d.ATExpiresIn = aTime
	d.RTExpiresIn = rTime
	d.Utime = now
	err = d.Get()
	if err != nil {
		err = d.Save()
	} else {
		err = d.Update()
	}
	if err != nil {
		return wxdata, err
	}
	wxdata.AccessToken  = access_token
	wxdata.ExpiresIn    = expires_in
	wxdata.RefreshToken = refresh_token
	return loginData(wxdata, userinfoResult), nil
}

//access_token登录
func loginByToken(access_token string) (wxLoginData, error) {
	wxdata := wxLoginData{}
	//var refresh_token string
	var expires_in int
	//通过id或access_token从数据库中获取refresh_token
	d := &data.WXLogin{AccessToken: access_token}
	err := d.GetByToken()
	if err != nil {
		return wxdata, err
	}
	//access_token是否有效
	err = WxLogin.VerifyAuth(d.OpenId, access_token)
	if err != nil {
		var now int64 = utils.Timestamp()
		//刷新access_token
		refreshResult, err := WxLogin.Refresh(d.RefreshToken)
		if err != nil {
			return wxdata, err
		}
		access_token = refreshResult.AccessToken
		expires_in = refreshResult.ExpiresIn
		//更新数据
		d.Set(access_token, now + int64(expires_in), now)
	}
	//获取个人信息
	userinfoResult, err := WxLogin.UserInfo(d.OpenId, access_token)
	if err != nil {
		return wxdata, err
	}
	wxdata.AccessToken  = access_token
	wxdata.ExpiresIn    = expires_in
	wxdata.RefreshToken = d.RefreshToken
	return loginData(wxdata, userinfoResult), nil
}

func weixinLogin(ctos *protocol.CWechatLogin, c interfacer.IConn) {
	stoc := &protocol.SWechatLogin{}
	var code string = ctos.GetCodeId()
	var token string = ctos.GetToken()
	glog.Infof("weixinLogin code:%s, token:%s", code, token)
	var member data.User
	//token登录
	if token != "" {
		wxdata, err := loginByToken(token)
		if err != nil {
			glog.Infof("weixinLogin err : %v", err)
			stoc.Error = proto.Uint32(uint32(protocol.Error_WechatLoingFailReAuth))
		} else {
			member, err = Login(wxdata, c)
			if err != nil {
				glog.Infof("weixinLogin err : %v", err)
				stoc.Error = proto.Uint32(uint32(protocol.Error_GetWechatUserInfoFail))
			}
			token = wxdata.AccessToken
		}
	}
	//code登录
	if code != "" {
		wxdata, err := loginByCode(code)
		if err != nil {
			glog.Infof("weixinLogin err : %v", err)
			stoc.Error = proto.Uint32(uint32(protocol.Error_WechatLoingFailReAuth))
		} else {
			member, err = Login(wxdata, c)
			if err != nil {
				glog.Infof("weixinLogin err : %v", err)
				stoc.Error = proto.Uint32(uint32(protocol.Error_GetWechatUserInfoFail))
			}
			token = wxdata.AccessToken
		}
	}
	if stoc.Error != nil {
		c.Send(stoc)
		time.AfterFunc(time.Millisecond*200, c.Close)
		return
	}
	logining(&member, c) //登录
	stoc.Userid = proto.String(member.Userid)
	stoc.Token = proto.String(token)
	c.Send(stoc)
}

func Login(wxdata wxLoginData, c interfacer.IConn) (data.User, error) {
	var member data.User
	err := member.GetByWechat(wxdata.OpenId)
	if err == nil {
		c.SetUserid(member.Userid)
		return member, nil
	}
	userid, err := data.GenUserID()
	if err != nil {
		return member, err
	}
	c.SetUserid(userid)
	ip := c.GetIPAddr()
	member = data.User{
		Wechat_uid:  wxdata.OpenId,
		Nickname:    wxdata.Nickname,
		Photo:       wxdata.HeadImagUrl,
		Sex:         uint32(wxdata.Sex),
		Userid:      userid,
		Create_ip:   ip,
		RoomCard:     20000,
		Create_time: uint32(time.Now().Unix()),
		Platform:    1,
	}
	err = member.Save()
	return member, err
}
