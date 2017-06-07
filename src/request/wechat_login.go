/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-01-23 10:22
 * Filename      : login.go
 * Description   : 微信登录
 * *******************************************************/
package request

import (
	"lib/socket"
	"crypto/tls"
	"data"
	"encoding/json"
	"errors"
	"fmt"
	"interfacer"
	"io/ioutil"
	"net/http"
	"protocol"
	"strconv"
	"time"

	"github.com/golang/glog"
	"code.google.com/p/goprotobuf/proto"
)

const targetUrl = "https://api.weixin.qq.com/sns/oauth2/access_token"
const userinfoUrl = "https://api.weixin.qq.com/sns/userinfo"
const refreshTockenUrl = "https://api.weixin.qq.com/sns/oauth2/refresh_token"

func init() {

	socket.Regist( &protocol.CWechatLogin{}, wechatLogin)
}

type WechatRefreshTokenRet struct {
	ErrCode       int    `json:"errcode"`
	Access_token  string `json:"access_token"`
	Expires_in    int    `json:"expires_in"`
	Refresh_token string `json:"refresh_token"`
	Openid        string `json:"openid"`
	Scope         string `json:"scope"`
}
type WechatLoginRet struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
}

type WechatUserInfoRet struct {
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

//?appid=APPID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN
func refreshToken(appid, token string) (*WechatRefreshTokenRet, error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	turl := refreshTockenUrl + fmt.Sprintf("?appid=%s&grant_type=refresh_token&refresh_token=%s", appid, token)
	req, err := http.NewRequest("GET", turl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		//glog.Infoln(string(data), appid, token)
		//解析返回的JSON数据
		var message WechatRefreshTokenRet
		err = json.Unmarshal(data, &message)
		if err == nil {
			if message.ErrCode == 0 {
				return &message, nil
			}
			err = errors.New("error code: " + strconv.Itoa(message.ErrCode))
		}
	}
	return nil, err
}
func getwechatUserinfo(token, openid string) (*WechatUserInfoRet, error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	turl := userinfoUrl + fmt.Sprintf("?access_token=%s&openid=%s", token, openid)
	req, err := http.NewRequest("GET", turl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		//glog.Infoln(err, string(data))
		//解析返回的JSON数据
		var message WechatUserInfoRet
		err = json.Unmarshal(data, &message)
		if err == nil {
			return &message, nil
		} else {
			return nil, err
		}
	}
	return nil, err
}

func getWechatAuth(appid, secret, code, grant_type string) (*WechatLoginRet, error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	turl := targetUrl + fmt.Sprintf("?appid=%s&secret=%s&code=%s&grant_type=%s", appid, secret, code, grant_type)
	req, err := http.NewRequest("GET", turl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	//	resp, err = client.PostForm(loginUrl, loginData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if data, err := ioutil.ReadAll(resp.Body); err == nil {
		//glog.Infoln("微信登录验证:", err, string(data))
		//解析返回的JSON数据
		var message WechatLoginRet
		err = json.Unmarshal(data, &message)
		//glog.Infoln("微信登录验证解码:", err, message)
		if err == nil {
			return &message, nil
		} else {
			return nil, err
		}
	}
	return nil, err
}

func wechatLogin(ctos *protocol.CWechatLogin, c interfacer.IConn) {
	stoc := &protocol.SWechatLogin{}
	var ret *WechatUserInfoRet
	//glog.Infoln(ctos.String())
	if ctos.GetToken() != "" {
		loginret, err := refreshToken(ctos.GetAppid(), ctos.GetToken())
		//glog.Infoln(loginret, err, ctos.GetAppid(), ctos.GetToken())
		if err != nil || loginret == nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_WechatLoingFailReAuth))
			c.Send(stoc)
			return

		}
		ret, err = getwechatUserinfo(loginret.Access_token, loginret.Openid)
		//glog.Infoln(ret, err, loginret.Access_token, loginret.Openid)
		if err != nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_GetWechatUserInfoFail))
			c.Send(stoc)
			return
		}
		stoc.Token = &loginret.Refresh_token

	} else {
		loginret, err := getWechatAuth(ctos.GetAppid(), ctos.GetSecret(), ctos.GetCodeId(), ctos.GetGrantType())
		//	glog.Infoln(loginret, err)
		if err != nil || loginret == nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_RegistError))
			c.Send(stoc)
			return
		}
		ret, err = getwechatUserinfo(loginret.AccessToken, loginret.OpenId)
		//	glog.Infoln(ret, err)
		if err != nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_GetWechatUserInfoFail))
			c.Send(stoc)
			return
		}
		stoc.Token = &loginret.RefreshToken
	}

	var member data.User
	// 该微信账号没有注册过
	if member.GetByWechat(ret.OpenId) != nil {
		userid, err := data.GenUserID()
		if err != nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_GetWechatUserInfoFail))
		} else {
			c.SetUserid(userid)
			ip := c.GetIPAddr()
			member = data.User{
				Wechat_uid:  ret.OpenId,
				Userid:      userid,
				Nickname:    ret.Nickname,
				Create_ip:   ip,
				RoomCard:     20000,
				Create_time: uint32(time.Now().Unix()),
				Photo:       ret.HeadImagUrl,
				Sex:         uint32(ret.Sex),
				Platform:    1,
			}
			err = member.Save()
			if err != nil {
				stoc.Error = proto.Uint32(uint32(protocol.Error_GetWechatUserInfoFail))
			}
		}
		// 微信账号已经登录过我们游戏
	} else {
		glog.Infoln("微信账号已经登录过我们游戏")
		c.SetUserid(member.Userid)
	}
	if stoc.Error != nil {
		//glog.Infoln(stoc)
		c.Send(stoc)
		time.AfterFunc(time.Millisecond*200, c.Close)
		return
	}
	logining(&member, c)
	stoc.Userid = proto.String(member.Userid)
	//glog.Infoln(*stoc)
	c.Send(stoc)

}
