//wx pay 微信支付
package request

import (
	"bytes"
	"data"
	"fmt"
	"wxpay"
	"io/ioutil"
	"net/http"
	"players"
	"strconv"

	"github.com/golang/glog"
	"config"
	"interfacer"
	"resource"
)

var Apppay *wxpay.AppTrans //微信支付

func WxPayInit() {
	notifyUrl := "https://"+config.Opts().Server_ip+config.Opts().Pay_port+config.Opts().Pay_wx_pattern
	cfg := &wxpay.WxConfig{
		AppId:         "wx987717ca7b85a17f",
		AppKey:        "N83KFU8P0Z8BN7DK39KFHJ8V73JBW7K9",
		MchId:         "1456128702",
		NotifyPattern: config.Opts().Pay_wx_pattern,
		NotifyUrl:     notifyUrl,
		PlaceOrderUrl: "https://api.mch.weixin.qq.com/pay/unifiedorder",
		QueryOrderUrl: "https://api.mch.weixin.qq.com/pay/orderquery",
		TradeType:     "APP",
	}
	appTrans, err := wxpay.NewAppTrans(cfg)
	if err != nil {
		panic(err)
	}
	Apppay = appTrans
	go Apppay.RecvNotify(wxRecvTrade) //goroutine
}

// 接收交易结果通知
func wxRecvTrade(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/plain;;charset=UTF-8")
	var buf bytes.Buffer
	if r.Method == "POST" {
		if result, err := ioutil.ReadAll(r.Body); err == nil {
			glog.Infof("tradeResult -> %s", string(result))
			tradeResult, err := wxpay.ParseTradeResult(result)
			glog.Infof("tradeResult -> %#v", tradeResult)
			if err == nil {
				go wxpayCallback(&tradeResult) //发货
			} else {
				glog.Errorf("trade result err: %v", err)
			}
		}
	}
	r.Body.Close()
	fmt.Fprintf(&buf, wxpay.TradeRespXml())
	//fmt.Fprintf(&buf, wxpay.TradeRespXmlFail())
	w.Write(buf.Bytes())
}

func wxpayCallback(t *wxpay.TradeResult) {
	err := Apppay.RecvVerify(t)
	if err != nil {
		glog.Errorf("recv verify %#v, err:, %v", t, err)
		return
	}
	//sign
	tradeRecord := new(data.TradeRecord)
	tradeRecord.Id = t.OrderId
	tradeRecord.Transid = t.TransactionId
	//tradeRecord := &data.TradeRecord{
	//	Id: t.OrderId,
	//	Transid: t.TransactionId,
	//}
	err = tradeRecord.Get()
	if err != nil {
		//订单不存在或其它
		glog.Errorf("not exist orderid %#v, err:, %v", t, err)
		return
	}
	if tradeRecord.Result == 0 {
		//重复发货
		glog.Errorf("repeat resp %#v, err:, %v", t, err)
		return
	}
	//更新记录
	tradeRecord.Transtime = t.TimeEnd
	tradeRecord.Currency = t.FeeType
	tradeRecord.Paytype = 403 //t.TradeType == "APP"
	money, err := strconv.Atoi(t.TotalFee)
	if err != nil {
		glog.Errorf("wxpay Callback: %#v, err: %v", t, err)
	}
	tradeRecord.Money = uint32(money) //转换为分
	tradeRecord.Result = 0 //交易成功
	// 离线状态
	//player := players.Get(t.OpenId) //TODO:优化
	player := players.Get(tradeRecord.Userid)
	if player == nil {
		tradeRecord.Result = 3 //发货中
		tradeRecord.SaveTradeOff()
	}
	//交易成功
	if player != nil {
		wxQuerySend(0, t.OrderId, player) //消息推送
		sendGoods(player, tradeRecord)
	}
	//update record
	err = tradeRecord.Update()
	if err != nil {
		glog.Errorf("tradeRecord:%#v, err:%v", tradeRecord, err)
	}
}


//发货
func sendGoods(player interfacer.IPlayer, t *data.TradeRecord) {
	propid, err := strconv.Atoi(t.Itemid)
	if err != nil {
		glog.Errorf("Send Goods: %v, err: %v", t, err)
		return
	}
	var count int32 = int32(t.Diamond)
	resource.ChangeRes(player, uint32(propid), count, data.RESTYPE6)
}

// 登录检测
func tradeOff(player interfacer.IPlayer) {
	list, err := data.GetTradeOff(player.GetUserid())
	if err == nil {
		for _, v := range list {
			sendGoods(player, v)
		}
		data.DelTradeOff(player.GetUserid())
	}
}

