package main

import (
	"time"

	"gorm.io/gorm"
)

// User has one `Account` (has one), many `Pets` (has many) and `Toys` (has many - polymorphic)
// He works in a Company (belongs to), he has a Manager (belongs to - single-table), and also managed a Team (has many - single-table)
// He speaks many languages (many to many) and has many friends (many to many - single-table)
// His pet also has one Toy (has one - polymorphic)



var TimeInterval = &struct {
	Begin, End string
}{
	Begin: "begin",
	End:   "end",
}

var ReadOnlyMapping = &struct {
	CreatedAt, DeletedAt string
}{
	CreatedAt: "created_at",
	DeletedAt: "deleted_at",
}

var RecordEnum = struct {
	Channel struct {
		Feishu, Wechat, Dingtalk, Mail string
	}
	State struct {
		Success, Failed string
	}
}{
	Channel: struct {
		Feishu, Wechat, Dingtalk, Mail string
	}{
		Feishu:   "feishu",
		Wechat:   "wechat",
		Dingtalk: "dingtalk",
		Mail:     "mail",
	},
	State: struct {
		Success, Failed string
	}{
		Success: "ok",
		Failed:  "error",
	},
}

var ro = ReadOnlyMapping

type Record struct {
	ID          string `json:"id" db:"id"`
	From        string `json:"from" db:"from"`
	To          string `json:"to" db:"to"`
	Sender      string `json:"sender" db:"sender"`
	SenderIP    string `json:"sender_ip" db:"sender_ip"`
	Receiver    string `json:"receiver" db:"receiver"`
	Title       string `json:"title" db:"title"`
	ContentType string `json:"content_type" db:"content_type"`
	Content     string `json:"content" db:"content"`
	AlertType   string `json:"channel" db:"alert_type"`
	EventType   string `json:"type" db:"event_type"`
	State       string `json:"state" db:"state"`
	StateMsg    string `json:"state_msg" db:"state_msg"`
	ReadOnly
}

type ReadOnly struct {
	CreatedAt time.Time      `json:"created_at" gorm:"created_at"` // 生成日期
	DeletedAt gorm.DeletedAt `json:"-" gorm:"deleted_at"`
}

func (m *ReadOnly) Create() {
	m.CreatedAt = time.Now()
}

func RecordTableName() string {
	return "alert_record"
}

func (Record) TableName() string {
	return RecordTableName()
}
