package app

import (
	"context"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func setToMemcache(ctx context.Context, value string) {
	log.Infof(ctx, "set into memcache")
	item := &memcache.Item{
		Key:        GAEEntityKind + "_" + "_" + HelloKey + "_Key",
		Value:      []byte(value),
		Expiration: MemcacheExprieTime,
	}
	memcache.Set(ctx, item)
}
