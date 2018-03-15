package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func getHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	ctx := appengine.NewContext(r)
	resultString := getFromDatastore(ctx)
	fmt.Fprintln(w, "Hello, "+resultString)
}

func getHandleWithDatastore(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	resultString := ""
	ctx := appengine.NewContext(r)

	if item, err := memcache.Get(ctx, GAEEntityKind+"_"+"_"+HelloKey+"_Key"); err == memcache.ErrCacheMiss {
		resultString = getFromDatastore(ctx)
		setToMemcache(ctx, resultString)
	} else if err != nil {
		log.Errorf(ctx, "%v", err)
	} else {
		log.Infof(ctx, "get from memcache, the vaule is %q", item.Value)
		resultString = string(item.Value)
	}

	fmt.Fprintln(w, "Hello, "+resultString)
}

func postHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}
	ctx := appengine.NewContext(r)
	oobj := OutputObject{}

	postForm := json.NewDecoder(r.Body)
	defer r.Body.Close()
	postForm.Decode(&oobj)

	setToDatastore(ctx, oobj)
	if IsUsingMemcache {
		setToMemcache(ctx, oobj.Key)
	}

	fmt.Fprintln(w, "Update Succeed")
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	switch status {
	case http.StatusNotFound:
		fmt.Fprint(w, "404 Not Found")
	case http.StatusMethodNotAllowed:
		fmt.Fprint(w, "405 Method Not Allow")
	default:
		fmt.Fprint(w, "Bad Request")
	}
}
