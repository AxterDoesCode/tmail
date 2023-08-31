package tmailcache

import "sync"

type Cache struct {
    MsgCacheMu sync.Mutex
	MsgCache map[string]MsgCacheEntry
    MsgCacheDisplay []MsgCacheEntry
}

type MsgCacheEntry struct {
	Id          string
	To          string
	From        string
	Subject     string
	ContentType string
	Body        string
}

func NewCache() Cache {
	return Cache{
		MsgCache: make(map[string]MsgCacheEntry),
        MsgCacheDisplay: []MsgCacheEntry{},
	}
}

func (c *Cache) AddToMessageCache(m *MsgCacheEntry) {
	c.MsgCache[m.Id] = *m
}

func (c *Cache) AddToMessageCacheDisplay(m *MsgCacheEntry) {
    c.MsgCacheDisplay = append(c.MsgCacheDisplay, *m)
}
