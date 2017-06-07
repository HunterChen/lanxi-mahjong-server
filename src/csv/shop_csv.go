package csv

import (
	"lib/csv"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"sync"
)

var shopMap map[uint32]ShopData
var shopLock sync.RWMutex
type ShopData struct {
	Id           uint32 `csv:"id"`          //
	PropId       uint32 `csv:"propid"`      // 兑换的物品,1=钻石,2=房卡
	Number       uint32 `csv:"number"`      // 兑换的数量
	Paymenttype  uint32 `csv:"paymenttype"` // 支付方式,1=RMB,2=钻石
	Price        uint32 `csv:"price"`       // 支付价格
}

func GetShop(id uint32) *ShopData {
	shopLock.RLock()
	defer shopLock.RUnlock()
	shop := shopMap[id]
	return &shop
}
func InitShop() {
	shopLock.Lock()
	defer shopLock.Unlock()
	f, err := os.Open("./csv/shop.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var shop []ShopData
	err = csv.Unmarshal(data, &shop)
	if err != nil {
		panic(err)
	}
	shopMap = make(map[uint32]ShopData)
	for _, v := range shop {
		shopMap[v.Id] = v
	}
	glog.Infoln("shop表：", len(shopMap))
}
