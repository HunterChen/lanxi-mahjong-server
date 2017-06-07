package wxpay

import (
	"testing"
	"encoding/json"
)

var Apppay *AppTrans //微信支付

func WxPayInit() {
	//host := data.Conf.Host
	//port := strconv.Itoa(data.Conf.PayPort)
	host := "119.29.24.17"
	port := "7005"
	pattern := "/mahjong/wxpay/notice"
	//pattern := data.Conf.PayWxPattern
	notifyUrl := "https://"+host+":"+port+pattern
	cfg := &WxConfig{
		AppId:         "wx987717ca7b85a17f",
		AppKey:        "N83KFU8P0Z8BN7DK39KFHJ8V73JBW7K9",
		MchId:         "1456128702",
		NotifyPattern: pattern,
		NotifyUrl:     notifyUrl,
		PlaceOrderUrl: "https://api.mch.weixin.qq.com/pay/unifiedorder",
		QueryOrderUrl: "https://api.mch.weixin.qq.com/pay/orderquery",
		TradeType:     "APP",
	}
	appTrans, err := NewAppTrans(cfg)
	if err != nil {
		panic(err)
	}
	Apppay = appTrans
	//go Apppay.RecvNotify(wxRecvTrade) //goroutine
}

func TestSubmit(t *testing.T) {
	WxPayInit()
	orderid := "SDGHUVNHiWb1159562"
	price := 600
	body := "购买12钻石"
	ip := "119.137.52.8"
	transid, err := Apppay.Submit(orderid, float64(price), body, ip)
	t.Log(err)
	t.Log(transid)
}

func TestReq(t *testing.T) {
	WxPayInit()
	orderid := "SDGHUVNHiWb1159562"
	payRequest := Apppay.NewPaymentRequest(orderid)
	t.Log(payRequest)
	retMap, err := ToMap(&payRequest)
	t.Log(err)
	t.Logf("retMap -> %#v", retMap)
	payReqStr := ToXmlString(retMap)
	t.Logf("payReqStr -> %s", payReqStr)
	retJson, err := ToJson(&payRequest)
	t.Log(err)
	t.Logf("retJson -> %s", string(retJson))
}

func TestJson(t *testing.T) {
	str := "{\"AppId\":\"wx987717ca7b85a17f\",\"PartnerId\":\"1456128702\",\"PrepayId\":\"wx201704011127404c2c8d2cd70314490910\",\"Package\":\"Sign=WXPay\",\"NonceStr\":\"90e10cf4fffb9b234b711f4fa5d8cce8\",\"Timestamp\":\"1491046061\",\"Sign\":\"CDEA6B79B09673636D24C7856A107CC8\"}"
	req := new(PaymentRequest)
	err := json.Unmarshal([]byte(str), req)
	t.Log(err)
	t.Logf("retJson -> %#v", req)
}

func TestTrade(t *testing.T) {
	WxPayInit()
	tradeResult := &TradeResult{
		ReturnCode: "SUCCESS",
		AppId: "wx987717ca7b85a17f",
		MchId: "1456128702",
		NonceStr:  "f15fb11fe6ee616c2c6879302ecba13b",
		Sign:  "C3377D79C3C4A4C2BFFDF22C1549CD6C",
		ResultCode: "SUCCESS",
		OpenId: "",
		IsSubscribe: "N",
		TradeType: "APP",
		BankType: "CMB_CREDIT",
		TotalFee: "600",
		FeeType: "CNY",
		CashFee: "600",
		TransactionId: "4008062001201704015410962069",
		OrderId: "XxZJhQN4iWb115956",
		Attach: "2",  
		TimeEnd: "20170401124954",
	}
	resultInMap := tradeResult.ToMap()
	t.Logf("resultInMap -> %#v", resultInMap)
	t.Log(resultInMap["sign"])
	wantSign := Sign(resultInMap, Apppay.Config.AppKey)
	t.Log(wantSign)
}

func TestCall(t *testing.T) {
	WxPayInit()
	xml := `<xml>
	<appid><![CDATA[wx987717ca7b85a17f]]></appid>
	<bank_type><![CDATA[CFT]]></bank_type>
	<cash_fee><![CDATA[600]]></cash_fee>
	<fee_type><![CDATA[CNY]]></fee_type>
	<is_subscribe><![CDATA[N]]></is_subscribe>
	<mch_id><![CDATA[1456128702]]></mch_id>
	<nonce_str><![CDATA[6d34d556070dda3c9e3b489fbf387854]]></nonce_str>
	<openid><![CDATA[oRALIw53XTeq1zi_MgvR_1vZ76QU]]></openid>
	<out_trade_no><![CDATA[FVpkePbCkWb16561]]></out_trade_no>
	<result_code><![CDATA[SUCCESS]]></result_code>
	<return_code><![CDATA[SUCCESS]]></return_code>
	<sign><![CDATA[0B8449A3DC86C720F65C0BFB2FA9A037]]></sign>
	<time_end><![CDATA[20170405110707]]></time_end>
	<total_fee>600</total_fee>
	<trade_type><![CDATA[APP]]></trade_type>
	<transaction_id><![CDATA[4010012001201704055921242348]]></transaction_id>
	</xml>`
	tradeResult, err := ParseTradeResult([]byte(xml))
	t.Log(err)
	t.Logf("tradeResult -> %#v", tradeResult)
	resultInMap := tradeResult.ToMap()
	t.Logf("resultInMap -> %#v", resultInMap)
	t.Log(resultInMap["sign"])
	wantSign := Sign(resultInMap, Apppay.Config.AppKey)
	t.Log(wantSign)
}
