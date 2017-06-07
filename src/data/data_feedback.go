package data

import (
	"lib/utils"
	"github.com/golang/glog"
)

type DataFeedback struct {
	Content    string //
	Createtime uint32 //
	Userid     string
	Status     uint32
	Kind       uint32
	Reply      string
	Replytime  uint32
	Replyer    string
	ImagePath  string
}

type DataNotice struct {
	List []*Notice
}

type Notice struct {
	Id      uint32 // time.Now().Unix()
	Type    uint32
	Title   string
	Content string
	CTime   uint32
	Expire  uint32
}

//
//func (this *DataFeedback) Get() error {
//	value, err := gossdb.C().Get(KEY_FEEDBACK + this.Userid + ":" + strconv.Itoa(int(this.Createtime)))
//	if err != nil {
//		return err
//	}
//	var data []byte = value
//	err = json.Unmarshal(data, this)
//	return err
//}
func (this *DataFeedback) Save() error {
	//this.Createtime = uint32(time.Now().Unix())
	//_, err := gossdb.C().Qpush_front(KEY_FEEDBACK, this)
	return nil
}

// notice
func (this *DataNotice) Get() error {

	return nil
	// return gossdb.C().GetObject(KEY_NOTICE, this)
}

func (this *DataNotice) Save() error {
	return nil
}

func (this *DataNotice) GetList() []*Notice {
	err := this.Get()
	if err != nil {
		glog.Infoln("GetList get error:", err)
	}
	now := uint32(utils.Timestamp())
	var list []*Notice
	for _, v := range this.List {
		if v.Expire < now {
			continue
		}
		list = append(list, v)
	}
	this.List = list
	err = this.Save()
	if err != nil {
		glog.Infoln("GetList save error:", err)
	}
	return this.List
}

func (this *DataNotice) Add(n *Notice) error {
	err := this.Get()
	if err != nil {
		glog.Infoln("Add get error:", err)
	}
	this.List = append(this.List, n)
	err = this.Save()
	return err
}
