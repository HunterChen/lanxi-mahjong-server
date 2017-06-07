package wxapi

import (
	"testing"
)

var WxLogin *Wxapi //微信登录

func WxLoginInit() {
	cfg := &WxapiConfig{
		AppId:           "wx1c3fa8244c190785",
		AppSecret:       "e25e0b8052a6cbe89be8921b96557728",
		AccessUrl:       "https://api.weixin.qq.com/sns/oauth2/access_token",
		RefreshUrl:      "https://api.weixin.qq.com/sns/oauth2/refresh_token",
		UserinfoUrl:     "https://api.weixin.qq.com/sns/userinfo",
		VerifyAccessUrl: "https://api.weixin.qq.com/sns/auth",
	}
	wx, err := NewWxapi(cfg)
	if err != nil {
		panic(err)
	}
	WxLogin = wx
}

func TestRefresh(t *testing.T) {
	WxLoginInit()
	openid := "o48K-wCLzOgvD7b2_kbllFcNHDmQ"
	access_token := "o2uq5NSy6ZqnsgshDSdlMA1bbZ-NQJAhS0TBH3rB9R9oyAv1SfrsYpNw3IM-8tbMZSzuFwHLG1uJmWVaLsqFlnEZiNpa1NvG5lJFHo3StyY"
	err := WxLogin.VerifyAuth(openid, access_token)
	t.Log(err)
	refresh_token := "o2uq5NSy6ZqnsgshDSdlMDRpcK0kb80dDwI9tlYxRhUYxXUdzhd7VbB3N8X_o8qLawEsSYdOZvhjel0SYU8MxprBjcdyrIWjl-4TBtLv2B0"
	refreshResult, err := WxLogin.Refresh(refresh_token)
	t.Log(err)
	t.Log(refreshResult)
	userinfoResult, err := WxLogin.UserInfo(openid, access_token)
	t.Log(err)
	t.Log(userinfoResult)
}

/*
appid:"wx1c3fa8244c190785" secret:"e25e0b8052a6cbe89be8921b96557728" code_id:"" grant_type:"authorization_code" token:"o2uq5NSy6ZqnsgshDSdlMDRpcK0kb80dDwI9tlYxRhUYxXUdzhd7VbB3N8X_o8qLawEsSYdOZvhjel0SYU8MxprBjcdyrIWjl-4TBtLv2B0" 
{"openid":"o48K-wCLzOgvD7b2_kbllFcNHDmQ","access_token":"o2uq5NSy6ZqnsgshDSdlMA1bbZ-NQJAhS0TBH3rB9R9oyAv1SfrsYpNw3IM-8tbMZSzuFwHLG1uJmWVaLsqFlnEZiNpa1NvG5lJFHo3StyY","expires_in":7200,"refresh_token":"o2uq5NSy6ZqnsgshDSdlMDRpcK0kb80dDwI9tlYxRhUYxXUdzhd7VbB3N8X_o8qLawEsSYdOZvhjel0SYU8MxprBjcdyrIWjl-4TBtLv2B0","scope":"snsapi_base,snsapi_userinfo,"}

<nil> {"openid":"o48K-wCLzOgvD7b2_kbllFcNHDmQ","nickname":"ASXCE","sex":1,"language":"zh_CN","city":"Wenzhou","province":"Zhejiang","country":"CN","headimgurl":"http:\/\/wx.qlogo.cn\/mmopen\/pDIUXgK1u9XiagggcVAN0pEibsiaxo0vZT9j69u8kIqgAYp5CdqXtoBuFbcyVvFicKGddxp9yHZmUTKKbEeBSqyiaYLaChMYHN2Vr\/0","privilege":[],"unionid":"otc1gwmEX8sbH1VCpGLiJMtzJp_Q"}
*/
