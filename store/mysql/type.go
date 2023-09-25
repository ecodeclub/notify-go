package mysql

import "time"

type Delivery struct {
	Id          int64     `xorm:"not null pk autoincr"`
	TemplateId  int64     `xorm:"INT"`
	Status      int       `xorm:"INT"`
	SendChannel int       `xorm:"comment('消息发送渠道 10.IM 20.Push 30.短信 40.Email 50.公众号') VARCHAR(255)"`
	MsgType     int       `xorm:"comment('10.通知类消息 20.营销类消息 30.验证码类消息') INT"`
	Proposer    string    `xorm:"comment('业务方') VARCHAR(255)"`
	Creator     string    `xorm:"VARCHAR(255)"`
	Updator     string    `xorm:"VARCHAR(255)"`
	IsDelted    int       `xorm:"INT"`
	Created     time.Time `xorm:"TIMESTAMP"`
	Updated     time.Time `xorm:"TIMESTAMP"`
}

type Target struct {
	Id           int64  `xorm:"not null pk INT"`
	TargetIdType string `xorm:"comment('接收目标id类型') VARCHAR(8)"`
	TargetId     string `xorm:"comment('接收目标id') VARCHAR(255)"`
	DeliveryId   int64  `xorm:"INT"`
	Status       int    `xorm:"INT"`
	MsgContent   string `xorm:"TEXT"`
}

type Template struct {
	Id         int64     `xorm:"not null pk INT"`
	Country    string    `xorm:"VARCHAR(255)"`
	Type       int       `xorm:"comment('sms|email|push') INT"`
	EnContent  string    `xorm:"TEXT"`
	ChsContent string    `xorm:"TEXT"`
	ChtContent string    `xorm:"TEXT"`
	Creator    string    `xorm:"VARCHAR(255)"`
	Updator    string    `xorm:"VARCHAR(255)"`
	IsDelted   int       `xorm:"INT"`
	Created    time.Time `xorm:"TIMESTAMP"`
	Updated    time.Time `xorm:"TIMESTAMP"`
}

func (Delivery) TableName() string {
	return "delivery"
}

func (Target) TableName() string {
	return "target"
}

func (Template) TableName() string {
	return "template"
}
