package cache

type Cache struct {
	Data map[string]any
}

func New() *Cache {
	return &Cache{
		Data: make(map[string]any),
	}
}

func (c *Cache) Set(uid string, data any) {
	c.Data[uid] = data
}

func (c *Cache) Get(uid string) (any, bool) {

	value, ok := c.Data[uid]

	if !ok {
		return nil, false
	}

	return value, true
}
