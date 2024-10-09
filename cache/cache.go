package cache

import (
	"fmt"
	"time"
)

type Cache interface {
	// 缓存数据
	Set(key string, val interface{}, timeout time.Duration)
	// 获取数据
	// 返回数据和数据是否存在
	Get(key string) (interface{}, bool)
}

// access_token 缓存 KEY
func (cli *Client) tokenCacheKey() string {
	return fmt.Sprintf("weapp.%s.access.token", cli.appid)
}
