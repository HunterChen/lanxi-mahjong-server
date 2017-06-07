/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2017-02-17 15:18
 * Filename      : wxpay.go
 * Description   : 微信支付
 * *******************************************************/
package request

import (
	"lib/socket"
	"lib/utils"
	"csv"
	"data"

	"interfacer"
	"players"
	"protocol"
	"strconv"
	"time"
	"wxpay"

	"github.com/golang/glog"
	"code.google.com/p/goprotobuf/proto"
)

func init() {

	socket.Regist(&protocol.CWxpayOrder{}, wxOrder)

	socket.Regist(&protocol.CWxpayQuery{}, wxQuery)
}

//微信支付查询
func wxQuery(ctos *protocol.CWxpayQuery, c interfacer.IConn) {
	stoc := &protocol.SWxpayQuery{Error: proto.Uint32(0)}
	var transid string = ctos.GetTransid()
	if transid == "" {
		stoc.Error = proto.Uint32(uint32(protocol.Error_IpayOrderFail))
	} else {
		queryResult, err := Apppay.Query(transid)
		if err != nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_IpayOrderFail))
		} else {
			if queryResult.ReturnCode == "SUCCESS" && queryResult.ResultCode == "SUCCESS" {
				stoc.Result = proto.Uint32(0)
				stoc.Orderid = proto.String(queryResult.OrderId)
			} else {
				stoc.Error = proto.Uint32(uint32(protocol.Error_IpayOrderFail))
			}
		}
	}
	c.Send(stoc)
}

//微信支付结果主动推送
func wxQuerySend(result uint32, orderid string, player interfacer.IPlayer) {
	stoc := &protocol.SWxpayQuery{Error: proto.Uint32(0)}
	stoc.Result = proto.Uint32(0)
	stoc.Orderid = proto.String(orderid)
	player.Send(stoc)
}

//微信支付下单
func wxOrder(ctos *protocol.CWxpayOrder, c interfacer.IConn) {
	stoc := &protocol.SWxpayOrder{Error: proto.Uint32(0)}
	var waresid uint32 = ctos.GetId()
	var body string = ctos.GetBody()
	var userid string = c.GetUserid()
	player:= players.Get(userid)
	var ipaddr uint32 = player.GetConn().GetIPAddr()
	var ip string = utils.InetTontoa(ipaddr).String()
	var parent string = player.GetBuild()
	for i := 0; i < 5; i++ { //失败重复尝试
		transid, orderid := wxpayOrder(waresid, userid, ip, parent, body)
		if transid == "" || orderid == "" {
			<-time.After(time.Duration(100) * time.Millisecond)
			glog.Infoln("wx order fail:", waresid, userid)
			stoc.Error = proto.Uint32(uint32(protocol.Error_IpayOrderFail))
		} else {
			payRequest := Apppay.NewPaymentRequest(transid)
			payReqJson, err := wxpay.ToJson(&payRequest)
			//retMap, err := wxpay.ToMap(&payRequest)
			if err != nil {
				glog.Infoln("wx order err:", waresid, userid, err)
				stoc.Error = proto.Uint32(uint32(protocol.Error_IpayOrderFail))
				break
			}
			//payReqStr := wxpay.ToXmlString(retMap)
			//stoc.Payreq = proto.String(payReqStr)
			stoc.Payreq = proto.String(string(payReqJson))
			//stoc.Orderid = proto.String(orderid)
			//stoc.Noncestr = proto.String(wxpay.NewNonceString())
			break
		}
	}
	stoc.Id = proto.Uint32(waresid)
	c.Send(stoc)
}

// 下单
func wxpayOrder(waresid uint32, userid, ip, parent, body string) (string, string) {
	d := csv.GetShop(waresid)
	if d.Paymenttype != 1 {
		glog.Errorf("wxpay waresid: %d", waresid)
		return "", ""
	}
	var diamond uint32 = d.Number
	var price uint32 = uint32(d.Price * 100) //转换为分
	var itemid string = strconv.FormatInt(int64(d.PropId), 10)
	var orderid string = data.GenCporderid(userid)
	transid, err := Apppay.Submit(orderid, float64(price), body, ip)
	if err != nil {
		glog.Errorf("wxpay waresid: %d, err: %v", waresid, err)
		return "", ""
	}
	//var ctime string = utils.Unix2Str(utils.Timestamp())
	//transid,下单记录
	t := &data.TradeRecord{
		Id: orderid,
		Transid: transid,
		Userid: userid,
		Itemid: itemid,
		Amount: "1",
		Diamond: diamond,
		Money: price,
		Ctime: time.Now().Unix(),
		Result: 2, //2=下单状态
		Clientip: ip,
		Parent: parent,
	}
	err = t.Save() //TODO:优化,未支付订单
	if err != nil {
		glog.Errorf("wxpay waresid: %d, err: %v", waresid, err)
		return "", ""
	}
	return transid, orderid
}
