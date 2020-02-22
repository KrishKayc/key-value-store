package db

import "sync"

var cwmux sync.Mutex
var crmux sync.Mutex

//For 1GB file with max capped values, cache takes 20000 keys, so this can be in memory and increase get() performance
type memcache struct {
	cache map[string]int64
}

func newMemCache(cache map[string]int64) memcache {
	return memcache{cache: cache}
}

func (c memcache) get(key string) int64 {
	var pos int64
	crmux.Lock()
	pos = c.cache[key]
	crmux.Unlock()

	return pos
}
func (c memcache) write(key string, pos int64) {

	cwmux.Lock()
	c.cache[key] = pos
	cwmux.Unlock()

}

func (c memcache) delete(key string) {

	cwmux.Lock()
	crmux.Lock()
	delete(c.cache, key)
	crmux.Unlock()
	cwmux.Unlock()

}
