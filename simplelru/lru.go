package simplelru

import (
	"container/list"
	"errors"
)

// EvictCallback 当缓存实体被删除时，调用此方法
type EvictCallback func(key, value interface{})

type LRU struct {
	size      int
	evictList *list.List
	items     map[interface{}]*list.Element
	onEvict   EvictCallback
}

type entity struct {
	key   interface{}
	value interface{}
}

func NewLRU(size int, onEvict EvictCallback) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}
	return &LRU{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   onEvict,
	}, nil
}

// 清空LRU缓存
func (c *LRU) Purge() {
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entity).value)
		}
		delete(c.items, k)
	}
	c.evictList.Init()
}

// 新增值到缓存中, 如果未发生淘汰，则返回false, 否则返回true
func (c *LRU) Add(key, value interface{}) (evicted bool) {
	// 1. 检查key是否已存在, 如果存在则更新value值, 并将该value移动到表头
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*entity).value = value
		return false
	}

	// 2. key不存在, 则执行将value添加到表头
	ent := &entity{key, value}
	entity := c.evictList.PushFront(ent)
	c.items[key] = entity

	// 3. 判断长度是否超过阈值, 如果超过则删除表尾元素
	evict := c.evictList.Len() > c.size
	if evict {
		c.removeOldest()
	}
	return evict
}

// 获取key对应的value, 如果命中缓存，则将value移动到表头
func (c *LRU) Get(key interface{}) (value interface{}, ok bool) {
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value.(*entity).value, true
	}
	return
}

// 判断key是否在缓存中, 不更新元素排列顺序
func (c *LRU) Contains(key interface{}) (ok bool) {
	_, ok := c.items[key]
	return ok
}

// 获取key对应的value, 如果命中缓存，则不会将value移动到表头
func (c *LRU) Peek(key interface{}) (value interface{}, ok bool) {
	if ent, ok := c.items[key]; ok {
		return ent.Value.(*entity).value, true
	}
	return
}

func (c *LRU) Remove(key interface{}) (present bool) {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}

func (c *LRU) RemoveOldest() (key, value interface{}, ok bool) {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
		entity := ent.Value.(*entity)
		return entity.key, entity.value, true
	}
	return nil, nil, false
}

func (c *LRU) GetOldest() (key, value interface{}, ok bool) {
	ent := c.evictList.Back()
	if ent != nil {
		entity := ent.Value.(*entity)
		return entity.key, entity.value, true
	}
	return nil, nil, false
}

func (c *LRU) Keys() []interface{} {
	var (
		kes = make([]interface{}, len(c.items))
		i   = 0
	)
	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = entity.Value.(*entity).key
		i++
	}
	return keys
}

func (c *LRU) Len() int {
	return c.evictList.Len()
}

func (c *LRU) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

func (c *LRU) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	entity := e.Value.(*entity)
	delete(c.items, entity.key)
	if c.onEvict != nil {
		c.onEvict(entity.key, entity, value)
	}
}
