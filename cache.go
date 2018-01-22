package gocache

import (
	"fmt"
)

var (
	_driver Driver
	//StorePrefix 存储空间前缀
	StorePrefix = ""
	//StoreSuffix 存储空间后缀
	StoreSuffix = "_"
	//DriverName 驱动名称
	DriverName string
)

var (
	//ErrFormat 错误的格式
	ErrFormat = fmt.Errorf("error format")
	//ErrNotFound 未找到
	ErrNotFound = fmt.Errorf("not found")
	//ErrNotPointer 给的值不是指针
	ErrNotPointer = fmt.Errorf("not pointer")
)

//Register 注册驱动
func Register(name string, mydriver Driver) {
	if name == "" {
		panic("name is null")
	}
	if mydriver == nil {
		panic("Register driver is nil")
	}
	DriverName = name
	_driver = mydriver
}

//Open 开启存储器
func Open(options interface{}) error {
	if options == nil {
		panic("driver does not init")
	}
	return _driver.Open(options)
}

//Close 关闭存储器
func Close() error {
	return _driver.Close()
}

//Store 返回一个存储器
func Store(storeName string) (storer Storer) {
	if _driver == nil {
		panic("driver is nil")
	}
	return _driver.Store(storeName)
}
