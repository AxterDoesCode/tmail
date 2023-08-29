package tmailcache

type Cache struct {
	MsgCache map[string]MsgCacheEntry
}

type MsgCacheEntry struct {
	Id      string
	To      string
	From    string
	Subject string
	Body    string
}

func NewCache() Cache {
	return Cache{
		MsgCache: make(map[string]MsgCacheEntry),
	}
}

func (c *Cache) AddToMessageCache(m MsgCacheEntry) {
	c.MsgCache[m.Id] = m
}
