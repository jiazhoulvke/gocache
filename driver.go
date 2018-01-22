package gocache

//Driver 缓存驱动
type Driver interface {
	//Open 获取存储器连接
	Open(options interface{}) error
	//Store 获取一个存储器实例
	Store(storeName string) Storer
	//Close 关闭存储器连接
	Close() error
}

//Storer 存储器
type Storer interface {
	//Set 存储值
	Set(key string, value interface{}) error
	//Remember 存储值，超时后自动删除
	Remember(key string, value interface{}, expireTime int) error
	//Get 获取值
	Get(key string, value interface{}) error
	//GetDefault 获取值，如果不存在则存储默认值
	GetDefault(key string, value interface{}, defaultValue interface{}) error
	//Delete 删除一个数据
	Delete(key string) error
	//String 获取存储器名称
	String() string
}
