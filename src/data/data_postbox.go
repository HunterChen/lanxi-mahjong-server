package data

import (
	"errors"
)

type DataPostbox struct {
	Id           uint32
	Title        string
	Appendixname string
	Appendix     []*WidgetData //
	Content      string
	Sender       string
	Receiver     string
	Createtime   uint32 //
	Read         bool   //
	Expire       uint32 //
	Kind         uint32 //
	Draw         bool   //
	Status       uint32
}

func (this *DataPostbox) ReadPost() error {
	err := this.Get()
	if err == nil {
		this.Status = 1
		err = this.Save()
	}
	return err
}

func (this *DataPostbox) Get() error {
	return nil
}

func (this *DataPostbox) Save() error {
	//this.Id = this.Receiver + this.Sender + strconv.Itoa(int(uint16(utils.RandInt64())))
	return nil
}

// Delete delete an email
func (this *DataPostbox) Delete() error {
	return nil
}

// Cleanup cleanup all emails
func (this *DataPostbox) Cleanup() error {
	return nil
}

// CleanupRead cleanup all readed emails
func (this *DataPostbox) CleanupRead() error {
	return nil
}

func (this *DataPostbox) ReadAll() ([]*DataPostbox, error) {
	list := make([]*DataPostbox, 0)

	return list, nil
}

func (this *DataPostbox) OpenAppendix() ([]*WidgetData, error) {
	err := this.Get()
	if err != nil {

	} else {
		if len(this.Appendix) > 0 {
			this.Draw = true
			err = this.Save()
		} else {
			err = errors.New("appendix is empty")
		}
	}
	return this.Appendix, err
}
