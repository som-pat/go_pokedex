package pokecache

import "time"

type Cache struct{
	cache 	map[string]cache_entry
}

type cache_entry struct{
	val 	 [] byte
	instance time.Time
}

func CreateCache(inter_time time.Duration) Cache {
	c:=  Cache{cache: make(map[string]cache_entry),}
	go c.LooPurge(inter_time)
	return c
}

func (c *Cache) Add(key string, val []byte){
	c.cache[key] = cache_entry{
		val : val,
		instance: time.Now().UTC(),
	}
}

func (c *Cache) Get(key string) ([]byte , bool){
	cval, ok := c.cache[key]
	return cval.val, ok
}

func (c *Cache) LooPurge(inter_time time.Duration){
	ticker := time.NewTicker(inter_time)	
	for range ticker.C{
		c.Purge(inter_time)
		 
		}
}


func (c *Cache) Purge(inter_time time.Duration) {
	interval := time.Now().UTC().Add(-inter_time)
	for k, v := range c.cache{
		if v.instance.Before(interval){
			delete(c.cache,k)
		}
	}
}