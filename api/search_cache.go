package api

import (
	"github.com/mangenotwork/search/utils/logger"
	"sync"
	"time"
)

// 使用定期清的缓存算法

var UrlCacheObj = NewUrlCache()

type UrlBody struct {
	Body       interface{}
	Count      int
	Url        string
	Expiration int64 // 过期时间
}

// Expired 判断数据项是否已经过期
func (udp *UrlBody) Expired() bool {
	if udp.Expiration == 0 {
		return false
	}
	return time.Now().Unix() > udp.Expiration
}

type UrlCache struct {
	defaultExpiration time.Duration
	items             map[string]*UrlBody // 缓存数据项存储在 map 中
	mu                sync.RWMutex        // 读写锁
	gcInterval        time.Duration       // 过期数据项清理周期
	stopGc            chan bool
}

func NewUrlCache() *UrlCache {
	c := &UrlCache{
		defaultExpiration: 10 * time.Second, //
		gcInterval:        10 * time.Second, // 10s清一次
		items:             map[string]*UrlBody{},
		stopGc:            make(chan bool),
	}
	// 开始启动过期清理 goroutine
	go c.gcLoop()
	return c
}

// gcLoop 过期缓存数据项清理
func (c *UrlCache) gcLoop() {
	logger.Info("run 过期缓存数据项清理")
	ticker := time.NewTicker(c.gcInterval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.stopGc:
			ticker.Stop()
			return
		}
	}
}

// DeleteExpired 删除过期数据项
func (c *UrlCache) DeleteExpired() {
	now := time.Now().Unix()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			c.mu.Lock()
			delete(c.items, k)
			c.mu.Unlock()
		}
	}
}

// Set 设置缓存数据项，如果数据项存在则覆盖
func (c *UrlCache) Set(urlStr string, body *UrlBody) {
	var e int64
	d := c.defaultExpiration
	if d > 0 {
		e = time.Now().Add(d).Unix()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	body.Expiration = e
	c.items[urlStr] = body
}

// Get 获取数据项
func (c *UrlCache) Get(k string) (*UrlBody, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	if item.Expired() {
		return nil, false
	}
	return item, true
}
