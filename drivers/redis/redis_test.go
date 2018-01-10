package redis

import (
	"testing"
	"time"

	"github.com/jiazhoulvke/gocache"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRedis(t *testing.T) {
	Convey("Redis Driver", t, func() {
		var err error
		err = gocache.Open(Options{
			Host:        "127.0.0.1",
			Port:        6379,
			IdleTimeout: 60,
		})
		So(err, ShouldBeNil)
		storeName := "testStore"
		keyName := "testKey"
		originValue := "abcdefg"
		var value string
		gocache.Store(storeName).Delete(keyName) //无论如何删除一次key，保证不会有残留数据
		So(storeName, ShouldEqual, gocache.Store(storeName).String())
		err = gocache.Store(storeName).Get(keyName, &value) //获取一个不存在的key，应该报错
		So(err, ShouldNotBeNil)
		err = gocache.Store(storeName).Set(keyName, originValue) //存储一个值，应该不报错
		So(err, ShouldBeNil)
		err = gocache.Store(storeName).Get(keyName, &value) //获取一个存在的值，应该不报错
		So(err, ShouldBeNil)
		So(originValue, ShouldEqual, value)            //获取的值应该和存储的值相等
		err = gocache.Store(storeName).Delete(keyName) //删除一个存在的值，应该不报错
		So(err, ShouldBeNil)
		err = gocache.Store(storeName).Get(keyName, &value) //获取一个已被删除的key，应该报错
		So(err, ShouldNotBeNil)
		err = gocache.Store(storeName).Remember(keyName, &value, 3) //缓存一个值，时间为3秒
		So(err, ShouldBeNil)
		err = gocache.Store(storeName).Get(keyName, &value) //获取一个存在的值，应该不报错
		So(err, ShouldBeNil)
		time.Sleep(time.Duration(4) * time.Second)          //睡眠4秒，等待key过期
		err = gocache.Store(storeName).Get(keyName, &value) //获取一个理论上应该已经过期的key，应该报错
		So(err, ShouldNotBeNil)
		err = gocache.Store(storeName).GetDefault(keyName, &value, "new_value") //获取一个key，由于key不存在，value应该等于赋予的默认值
		So(err, ShouldBeNil)
		So(value, ShouldEqual, "new_value")
		err = gocache.Store(storeName).GetDefault(keyName, &value, "new_value2") //由于key已经存在，所以应该返回key值
		So(err, ShouldBeNil)
		So(value, ShouldEqual, "new_value")
	})
}
