package job

import (
	"ConDetect/backend/global"
	"time"
)

type Cache struct{}

func NewCacheJob() *Cache {
	return &Cache{}
}

func (c *Cache) Run() {
	global.LOG.Info("run cache gc start ...")
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
	again:
		err := global.CacheDb.RunValueLogGC(0.7)
		if err == nil {
			goto again
		}
	}
	global.LOG.Info("run cache gc end ...")
}
