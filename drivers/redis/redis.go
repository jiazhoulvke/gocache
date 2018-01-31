package redis

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jiazhoulvke/gocache"
	json "github.com/json-iterator/go"
)

var (
	//ErrOptionsFormat 参数格式错误
	ErrOptionsFormat = fmt.Errorf("options format error")
)

func init() {
	gocache.Register("redis", &Driver{})
}

//Driver redis驱动
type Driver struct {
	pool *redis.Pool
}

//Options options
type Options struct {
	Host        string
	Port        int
	IdleTimeout int
}

//Open 打开连接
func (d *Driver) Open(options interface{}) error {
	if options == nil {
		return fmt.Errorf("options is nil")
	}
	opts, ok := options.(Options)
	if !ok {
		return ErrOptionsFormat
	}
	if opts.IdleTimeout == 0 {
		opts.IdleTimeout = 1
	}
	d.pool = &redis.Pool{
		MaxIdle:     runtime.NumCPU(),
		IdleTimeout: time.Duration(opts.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			var c redis.Conn
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", opts.Host, opts.Port))
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

//Close 关闭连接
func (d *Driver) Close() error {
	return d.pool.Close()
}

//Store 返回一个存储器
func (d *Driver) Store(storeName string) gocache.Storer {
	store := &Store{
		Name:   storeName,
		Prefix: gocache.StorePrefix + storeName + gocache.StoreSuffix,
		driver: d,
	}
	return store
}

//Store store
type Store struct {
	Name   string
	Prefix string
	driver *Driver
}

//Key 获取真实的键名
func (s *Store) Key(key string) string {
	return s.Prefix + key
}

//Delete 删除键值对
func (s *Store) Delete(key string) error {
	c := s.driver.pool.Get()
	_, err := c.Do("DEL", s.Key(key))
	return err
}

//Get 获取值
func (s *Store) Get(key string, obj interface{}) error {
	c := s.driver.pool.Get()
	reply, err := c.Do("GET", s.Key(key))
	if err != nil {
		return err
	}
	if reply == nil {
		return gocache.ErrNotFound
	}
	if content, ok := reply.([]byte); ok {
		return json.Unmarshal(content, obj)
	}
	return gocache.ErrFormat
}

//GetDefault 获取值，如果不存在则存储默认值
func (s *Store) GetDefault(key string, obj interface{}, defaultValue interface{}) error {
	c := s.driver.pool.Get()
	reply, err := c.Do("GET", s.Key(key))
	if err != nil || reply == nil {
		err = s.Set(key, defaultValue)
		if err != nil {
			return err
		}
		v := reflect.ValueOf(obj)
		if !v.Elem().CanSet() {
			return gocache.ErrNotPointer
		}
		v.Elem().Set(reflect.ValueOf(defaultValue))
		return nil
	}
	if content, ok := reply.([]byte); ok {
		return json.Unmarshal(content, obj)
	}
	return gocache.ErrFormat
}

//Set 设置值
func (s *Store) Set(key string, value interface{}) error {
	return s.Remember(key, value, 0)
}

//Remember 存储值，超时后自动删除
func (s *Store) Remember(key string, value interface{}, expireTime int) error {
	c := s.driver.pool.Get()
	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.Do("SET", s.Key(key), bs)
	if expireTime > 0 {
		_, err := c.Do("EXPIRE", s.Key(key), expireTime)
		if err != nil {
			return err
		}
	}
	return err
}

//String string
func (s *Store) String() string {
	return s.Name
}
