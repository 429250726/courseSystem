package cache

import "sync"
import "courseSys/cache/lru"

type cache struct{
	mu sync.Mutex
	lru *lru.Cache
	cacheBytes int64
}

func (c *cache) Lock(){
	c.mu.Lock()
}

func (c *cache) Unlock(){
	c.mu.Unlock()
}

func (c *cache) add(key string,value ByteView){
	c.Lock()
	defer c.Unlock()
	if c.lru==nil{
		c.lru=lru.NewCache(c.cacheBytes,nil)
	}
	c.lru.Add(key,value)
}

func (c *cache) get(key string) (value ByteView,ok bool){
	c.Lock()
	defer c.Unlock()
	if c.lru==nil{
		return
	}
	if v,ok:=c.lru.Get(key);ok{
		return v.(ByteView),ok
	}
	return
}




