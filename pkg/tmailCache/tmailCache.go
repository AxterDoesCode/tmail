package tmailcache

type Cache struct {
    MsgCache map[string]MsgCacheEntry
}

type MsgCacheEntry struct {
    To string
    From string
    Subject string
    Body string
}

func NewCache () Cache {
    return Cache {
        MsgCache: make(map[string]MsgCacheEntry),
    }
}
