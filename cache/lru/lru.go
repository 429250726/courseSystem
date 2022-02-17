package lru

import "container/list"

type Cache struct{
	maxBytes int64
	nBytes int64
	ll *list.List
	cache map[string]*list.Element
	OnEvicted func(key string,value Valuer)
}

func NewCache(maxBytes int64,OnEvicted func(string,Valuer))*Cache{
	return &Cache{
		maxBytes: maxBytes,
		nBytes: 0,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

func (c *Cache) Get(key string) (value Valuer,ok bool){
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv:=ele.Value.(*keyValue)
		return kv.value,true
	}
	return
}

func (c *Cache) RemoveOldest(){
	ele:=c.ll.Back()
	if ele!=nil{
		c.ll.Remove(ele)
		kv:=ele.Value.(*keyValue)
		delete(c.cache,kv.key)
		c.nBytes-=int64(len(kv.key)+kv.value.Len())
		if c.OnEvicted!=nil{
			c.OnEvicted(kv.key,kv.value)
		}
	}
}

func (c *Cache) Add(key string,value Valuer){
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv:=ele.Value.(*keyValue)
		c.nBytes+=int64(value.Len()-kv.value.Len())
		kv.value=value
	}else{
		ele:=c.ll.PushFront(&keyValue{key,value})
		c.cache[key]=ele
		c.nBytes+=int64(len(key)+value.Len())
	}
	for c.maxBytes!=0 && c.maxBytes<c.nBytes{
		c.RemoveOldest()
	}
}













