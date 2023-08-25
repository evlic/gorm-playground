package main

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GORM_REPO: https://github.com/go-gorm/gorm.git
// GORM_BRANCH: master
// TEST_DRIVERS: sqlite, mysql, postgres, sqlserver

func TestGORM(t *testing.T) {
	TestBySQLite(t)
}

const (
	dsn = ":memory:"
	//logEnable = false
	//logEnable = true
)

var (
	c = context.TODO()
	d *dao
	r *Record

	stateEnum   = RecordEnum.State
	stateAll    = []string{stateEnum.Success, stateEnum.Failed}
	channelEnum = RecordEnum.Channel
	channelAll  = []string{
		channelEnum.Dingtalk,
		channelEnum.Feishu,
		channelEnum.Wechat,
		channelEnum.Mail,
	}
)

func TestMain(m *testing.M) {
	testDB, err := gorm.Open(sqlite.Open(dsn))

	if err != nil {
		panic(err)
	}

	//testDBLog(testDB, logEnable)
    testDBLogOn = func () {
        testDBLog(testDB, true)
    }
    testDBLogOff = func () {
        testDBLog(testDB, false)
    }

	d = &dao{testDB}
	if err = testDB.AutoMigrate(r); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func getRecord(i int) *Record {
	rod := &Record{
		ID:          uuid.New().String(),
		To:          uuid.New().String(),
		Sender:      uuid.New().String(),
		Receiver:    uuid.New().String(),
		Title:       uuid.New().String(),
		ContentType: uuid.New().String(),
		Content:     uuid.New().String(),
		AlertType:   channelAll[i%len(channelAll)],
		EventType:   uuid.New().String(),
		SenderIP:    uuid.New().String(),
		State:       stateAll[i%len(stateAll)],
		StateMsg:    uuid.New().String(),
	}
	rod.Create()
	return rod
}

const (
	n = 100
)

func TestBySQLite(t *testing.T) {
    //testDBLogOn()

	Convey("insert", t, func() {
       //testDBLogOff()
		for i := 0; i < n; i++ {
			r = getRecord(i)
			err := d.CreateRecord(c, r)
			So(err, ShouldBeNil)
		}
	})
    testDBLogOn()
	query := map[string]queryParam{
		// "with-time":    getParamWithTimeLimit(),
		// "filter-state": getParamWithState(),
		// "sort-time":    getParamWithSortByTime(),
		"search-x":     getParamWithSearch(),
	}
	Convey("query", t, func() {
		for key, v := range query {
			Convey("Testing " + key, func() {
				cnt, err := d.CountRecord(c, v.filter, v.or, v.time)
				So(err, ShouldBeNil)

				list, err := d.ListRecord(c, v.filter, v.or, v.time, v.sort, v.page, v.size)
                t.Log("Testing " + key, cnt, len(list), v.expectation)
				So(err, ShouldBeNil)
				So(cnt, ShouldEqual, len(list))
				So(cnt, ShouldEqual, v.expectation)
			})
		}
	})
}

var (
	testDBLogOn  = func() {}
	testDBLogOff = func() {}
)

func testDBLog(testDB *gorm.DB, enable bool) {
	if enable {
		testDB.Logger = logger.Default.LogMode(logger.Info)
		return
	}
	testDB.Logger = logger.Default.LogMode(logger.Silent)
}


type queryParam struct {
	filter map[string]any
	or     map[string]any
	time   map[string]any
	sort   string
	page   int
	size   int

	expectation int
}

func getParamWithTimeLimit() queryParam {
	now := time.Now()
	return queryParam{
		time: map[string]any{
			TimeInterval.Begin: now.Add(time.Hour * -24),
			TimeInterval.End:   now,
		},
		expectation: n,
	}
}

func getParamWithSortByTime() queryParam {
	//now := time.Now()
	return queryParam{
		sort:        "+" + ro.CreatedAt,
		expectation: n,
	}

}

func getParamWithSearch() queryParam {
	//now := time.Now()
	return queryParam{
		or: map[string]any{
			"id": "%x%",
			"\"from\"": "%x%",
		},
		expectation: 0,
	}
}

func getParamWithState() queryParam {
	//now := time.Now()

	return queryParam{
		filter: map[string]any{
			"state": []string{stateEnum.Failed},
		},
		expectation: n >> 1,
	}
}
