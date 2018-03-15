package app

import (
	"net/http"
	"time"

	"google.golang.org/appengine"
)

// GAEEntityKind is for GAE Entity kind
// HelloKey is for GAE Entity string id
// IsUsingMemcache is using Memcache or not
// MemcacheExprieTime is exprire time sec
const (
	GAEEntityKind      = "<DATASTORE_KIND>"
	HelloKey           = "<SET_YOUR_KEY>"
	IsUsingMemcache    = false
	MemcacheExprieTime = time.Duration(3600) * time.Second
)

func init() {
	if IsUsingMemcache {
		http.HandleFunc("/get", getHandleWithDatastore)
	} else {
		http.HandleFunc("/get", getHandle)
	}
	http.HandleFunc("/post", postHandle)
	appengine.Main()
}
