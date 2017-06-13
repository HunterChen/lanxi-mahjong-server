package data

import (
	"errors"
	"time"
)

type DataUserActive struct {
	Userid string //用户账号id
	IP     uint32
	Time   uint32 // 时间戳
	Action uint32 // 1:上线，2：下线
	Device string // 设备型号
}

// 记录登陆时间,没有该玩家数据则插入
func (this *DataUserActive) Login() error {
	t := time.Now().Unix()
	this.Time = uint32(t)
	this.Action = 1
	var err error
	if this.Userid != "" {
	} else {
		err = errors.New("user id is empty!")
	}
	return err
}

// 记录退出时间，并累积在线时长
func (this *DataUserActive) Logout() error {
	t := time.Now().Unix()
	this.Time = uint32(t)
	this.Action = 2

	var err error
	if this.Userid != "" {
	} else {
		err = errors.New("user id is empty!")
	}
	return err
}
