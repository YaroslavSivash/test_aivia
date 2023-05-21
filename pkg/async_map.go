package async_map

import "sync"

type Counters struct {
	mx sync.Mutex
	m  map[string]string
}

func NewCounters() *Counters {
	return &Counters{
		m: make(map[string]string),
	}
}

func (c *Counters) Load(key string) (string, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	val, ok := c.m[key]
	return val, ok
}

func (c *Counters) Store(key string, value string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[key] = value
}

func (c *Counters) LoadAll() map[string]string {
	c.mx.Lock()
	defer c.mx.Unlock()
	return c.m
}
