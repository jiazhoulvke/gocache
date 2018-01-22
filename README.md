# gocache #

golang的缓存库。目前仅支持redis，但可以通过自己编写driver进行扩展。

```go
package main

import (
	"fmt"

	"github.com/jiazhoulvke/gocache"
	"github.com/jiazhoulvke/gocache/drivers/redis"
)

func main() {
	if err := gocache.Open(redis.Options{
		Host:        "127.0.0.1",
		Port:        6379,
		IdleTimeout: 60,
	}); err != nil {
		panic(err)
	}
	if err := gocache.Store("test").Set("foo", "bar"); err != nil {
		panic(err)
	}
	var s string
	if err := gocache.Store("test").Get("foo", &s); err != nil {
		panic(err)
	} else {
		fmt.Println(s)
	}
}
```
