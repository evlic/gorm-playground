package main

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func main() {
	fmt.Println("vim-go")
}

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

type dao struct {
	orm *gorm.DB
}

const (
	QueryListBufferSize = 8
	QueryListMaximum = 10000
)

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

func (d *dao) ListRecord(c context.Context, inFilter, search, inTime map[string]any, sort string, limit, offset int) (res []*Record, err error) {
	tx := d.orm.WithContext(c).Table(RecordTableName())
	res = make([]*Record, 0, QueryListBufferSize)

	if len(sort) > 4 {
		tx = tx.Order(sort)
	}

	if inTime != nil {
		interval := TimeInterval
		tx = tx.Where(
			"date("+ro.CreatedAt+") between date(?) and date(?)",
			inTime[interval.Begin],
			inTime[interval.End],
		)
	}

	if search != nil {
		txOr := tx
		for filed, v := range search {
			txOr = txOr.Or(filed+" like ?", v)
		}
		tx.Where(txOr)
	}

	tx = tx.Where(inFilter)

	var total int64
	err = tx.Count(&total).Error
	if limit > 0 && offset > -1 {
		tx = tx.Limit(limit).Offset(offset)
	} else {
		tx = tx.Limit(QueryListMaximum)
	}

	err = tx.Find(&res).Error
	return
}

func (d *dao) CountRecord(c context.Context, filter, or, time map[string]any) (total int64, err error) {
	tx := d.orm.WithContext(c).Table(RecordTableName())
	if time != nil {
		interval := TimeInterval
		tx = tx.Where(
			"date("+ro.CreatedAt+") between date(?) and date(?)",
			time[interval.Begin],
			time[interval.End],
		)
	}

	filedCnt := 0
	for filed, v := range or {
		if filedCnt == 0 {
			tx = tx.Where(filed+" like ?", v)
			filedCnt++
			continue
		}
		tx = tx.Or(filed+" like ?", v)
		filedCnt++
	}

	tx = tx.Where(filter).Where("deleted_at is null")

	err = tx.Count(&total).Error
	return
}

func (d *dao) CreateRecord(c context.Context, r *Record) (err error) {
	err = d.orm.WithContext(c).Create(r).Error
	return
}
