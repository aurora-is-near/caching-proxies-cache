package cache

import (
	"context"
	"sync"
	"time"
)

type Element struct {
	Blob        []byte
	WasCachedAt time.Time
}

type Key struct {
	PreviousHashId string
	ShardId        string
}

type Cache struct {
	lock *sync.Mutex
	mp   map[Key]*Element
	ttl  time.Duration
}

func (c *Cache) Get(previousHashId string, shardId string) *Element {
	c.lock.Lock()
	defer c.lock.Unlock()

	entry := c.mp[Key{
		PreviousHashId: previousHashId,
		ShardId:        shardId,
	}]

	if entry == nil {
		return nil
	}

	if time.Since(entry.WasCachedAt) > c.ttl {
		// Not evicted. Yet.
		return nil
	}

	return entry
}

func (c *Cache) SetIfBlobIsBigger(previousHashId string, shardId string, blob []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	key := Key{
		PreviousHashId: previousHashId,
		ShardId:        shardId,
	}

	previousEntry := c.mp[key]
	if previousEntry != nil {
		if len(previousEntry.Blob) > len(blob) {
			return
		}
	}

	c.mp[key] = &Element{
		Blob:        blob,
		WasCachedAt: time.Now(),
	}
}

func (c *Cache) runEvictionWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			c.evict()
		}
	}
}

func (c *Cache) evict() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for k, v := range c.mp {
		if time.Since(v.WasCachedAt) > c.ttl {
			delete(c.mp, k)
		}
	}
}

func NewCache(ctx context.Context, duration time.Duration) *Cache {
	c := &Cache{
		lock: &sync.Mutex{},
		mp:   map[Key]*Element{},
		ttl:  duration,
	}

	go c.runEvictionWorker(ctx)
	return c
}
