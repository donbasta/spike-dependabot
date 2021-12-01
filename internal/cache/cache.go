package cache

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var c *cache.Cache
var once sync.Once

func ProvideCache() *cache.Cache {
	once.Do(func() {
		c = cache.New(5*time.Minute, 10*time.Minute)
	})
	return c
}
