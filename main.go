package main

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/playground/db"
)

func main() {
	fmt.Println("vim-go")
}

type dao struct {
	orm *gorm.DB
}


func (d *dao) ListRecord(c context.Context, inFilter, search, inTime map[string]any, sort string, limit, offset int) (res []*Record, err error) {
    tx := d.orm.WithContext(c).Table(RecordTableName())
    res = make([]*Record, 0, db.QueryListBufferSize)

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
        tx = tx.Limit(db.QueryListMaximum)
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

    //filedCnt := 0
    //for filed, v := range or {
    //	if filedCnt == 0 {
    //		tx = tx.Where(filed+" like ?", v)
    //		filedCnt++
    //		continue
    //	}
    //	tx = tx.Or(filed+" like ?", v)
    //	filedCnt++
    //}
    
    if or != nil {
        txOr := tx
        for filed, v := range or {
            txOr = txOr.Or(filed+" like ?", v)
        }
        tx.Where(txOr)
    }

    tx = tx.Where(filter).Where("deleted_at is null")

    err = tx.Count(&total).Error
    return
}

func (d *dao) CreateRecord(c context.Context, r *Record) (err error) {
    err = d.orm.WithContext(c).Create(r).Error
    return
}
