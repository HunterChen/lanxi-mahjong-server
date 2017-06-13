package room

import (
	"game/algorithm"
	"time"
	"math/rand"
	"game/resource"
	"game/data"
	"lib/utils"
	"github.com/golang/glog"
)

//发牌
func (t *Desk) deal() {
	for s, p := range t.players {
		var cards []byte = t.getHandCards(s)
		if t.dealer == s {
			//庄家提示处理
			v:= t.DrawDetect(0, cards, []uint32{}, []uint32{}, []uint32{}, t.luckyCard,t.dealer)
			v|=algorithm.DetectKong(cards,[]uint32{}, t.luckyCard)
			if v > 0 {
				t.opt[s] = v //设置操作状态值
			}
			t.draw = cards[len(cards)-1] //庄家最后一张默认为摸牌
			//庄家消息
			msg := res_zhuangDeal(v, t.dice, cards, uint32(t.luckyCard))
			p.Send(msg)
		} else {
			//闲家消息
			msg := res_deal(0, t.dice, cards, uint32(t.luckyCard))
			p.Send(msg)
		}
	}
}

//洗牌
func (t *Desk) shuffle() {
	rand.Seed(time.Now().UnixNano())
	d := make([]byte, algorithm.TOTAL, algorithm.TOTAL)
	copy(d, algorithm.CARDS)
	//测试暂时去掉洗牌
	for i := range d {
		j := rand.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
	t.cards = d
}

//开始游戏
func (t *Desk) gameStart() {
	glog.Infof("gameStart -> %d, seat -> %d", t.id, t.seat)
	t.gameStartInit() //初始化

	// 牌局开始扣除房主的房卡
	for _, p := range t.players {
		if t.round == 0 && t.data.Cid == p.GetUserid(){
			resource.ChangeRes(p, resource.ROOM_CARD, -1 * int32(t.data.Cost), data.RESTYPE4)
		}
	}

	//打骰(两个骰子)
	dice1 := uint32(utils.RandInt32N(5) + 1)
	dice2 := uint32(utils.RandInt32N(5) + 1)
	t.dice = (dice1 << 16) + dice2 //TODO:优化

	if len(t.cheatLeftCards) > 0 {
		t.dealer = 1
		for s, _ := range t.players {
			t.handCards[s] = make([]byte, len(t.cheatHandCards[s-1]))
			copy(t.handCards[s], t.cheatHandCards[s-1])
		}
		// 需要copy，和上局结束后的牌墙重叠了
		t.cards = make([]byte, len(t.cheatLeftCards))
		copy(t.cards, t.cheatLeftCards)

		msg := res_dealer(t.dealer)
		t.broadcast(msg) //打庄消息通知
		t.luckyCard = t.cheatwildcard
	} else {
		if t.dealer == 0 {
			t.dealer = uint32(utils.RandInt32N(4) + 1)
		}

		msg := res_dealer(t.dealer)
		t.broadcast(msg) //打庄消息通知
		t.shuffle()      //洗牌
		// 洗完牌再产生财神牌
		index := int(utils.RandInt32N(int32(len(t.cards))))
		t.luckyCard = t.cards[index]
		t.cards = append(t.cards[:index], t.cards[index+1:]...)

		for s, _ := range t.players {
			var hand int = int(algorithm.HAND)
			if s == t.dealer { //判断庄家发14张牌
				hand += 1
			}
			cards := make([]byte, hand, hand)
			tmp := t.cards[:hand]
			copy(cards, tmp)
			t.handCards[s] = cards
			t.cards = t.cards[hand:]
		}
	}

	//第一个操作为庄家
	t.seat = t.dealer
	t.deal() //发牌
}

//初始化
func (t *Desk) gameStartInit() {
	t.state = true //设置房间状态
	if t.closeCh == nil {
		t.closeCh = make(chan bool, 1)
	}
	t.operateInit()
	t.skip = make(map[uint32]uint32)
	t.outCards = make(map[uint32][]byte)    //海底牌
	t.pongCards = make(map[uint32][]uint32) //碰牌
	t.kongCards = make(map[uint32][]uint32) //杠牌
	t.chowCards = make(map[uint32][]uint32) //吃牌(8bit-8-8)
	t.handCards = make(map[uint32][]byte)   //手牌
}
